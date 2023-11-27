// Package file is a file source. Expected format is json
package file

import (
	"io"
	"os"

	"objgo/team/core/config/source"

	"github.com/fsnotify/fsnotify"
)

type file struct {
	path string
	opts source.Options
}

var (
	DefaultPath = "config.json"
)

func (f *file) Read() (*source.ChangeSet, error) {
	fh, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()
	b, err := io.ReadAll(fh)
	if err != nil {
		return nil, err
	}
	info, err := fh.Stat()
	if err != nil {
		return nil, err
	}

	cs := &source.ChangeSet{
		Format:    format(f.path, f.opts.Encoder),
		Source:    f.String(),
		Timestamp: info.ModTime(),
		Data:      b,
	}
	cs.Checksum = cs.Sum()

	return cs, nil
}

func (f *file) String() string {
	return "file"
}

func (f *file) Watch() (source.Watcher, error) {
	if _, err := os.Stat(f.path); err != nil {
		return nil, err
	}
	return newWatcher(f)
}

func (f *file) Write(cs *source.ChangeSet) error {
	return nil
}

func NewSource(opts ...source.Option) source.Source {
	options := source.NewOptions(opts...)
	path := DefaultPath
	f, ok := options.Context.Value(filePathKey{}).(string)
	if ok {
		path = f
	}
	return &file{opts: options, path: path}
}

type watcher struct {
	f *file

	fw   *fsnotify.Watcher
	exit chan bool
}

func newWatcher(f *file) (source.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	err = fw.Add(f.path)
	if err != nil {
		return nil, err
	}

	return &watcher{
		f:    f,
		fw:   fw,
		exit: make(chan bool),
	}, nil
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	// is it closed?
	select {
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	default:
	}

	// try get the event
	select {
	case event, _ := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			// check existence of file, and add watch again
			_, err := os.Stat(event.Name)
			if err == nil || os.IsExist(err) {
				err := w.fw.Add(event.Name)
				if err != nil {
					return nil, err
				}
			}
		}

		c, err := w.f.Read()
		if err != nil {
			return nil, err
		}
		return c, nil
	case err := <-w.fw.Errors:
		return nil, err
	case <-w.exit:
		return nil, source.ErrWatcherStopped
	}
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}
