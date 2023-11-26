package api

import (
	"objgo/common/database"
	"objgo/team/core/config/source/file"
	"objgo/team/core/sdk/config"

	"github.com/spf13/cobra"
)

var (
	configYml string
	apiCheck  bool
	StartCmd  = &cobra.Command{
		Use:          "server",
		Short:        "Start API server",
		Example:      "go-admin server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().BoolVarP(&apiCheck, "api", "a", false, "Start server with check api data")

	//注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	// AppRouters = append(AppRouters, router.InitRouter)
}

func setup() {
	// 注入配置扩展项
	// config.ExtendConfig = &ext.ExtConfig
	//1. 读取配置
	config.Setup(
		file.NewSource(file.WithPath(configYml)),
		database.Setup,
		storage.Setup,
	)
	//注册监听函数
	// queue := sdk.Runtime.GetMemoryQueue("")
	// queue.Register(global.LoginLog, models.SaveLoginLog)
	// queue.Register(global.OperateLog, models.SaveOperaLog)
	// queue.Register(global.ApiCheck, models.SaveSysApi)
	// go queue.Run()

	// usageStr := `starting api server...`
	// log.Println(usageStr)
}

func run() error {

	return nil
}
