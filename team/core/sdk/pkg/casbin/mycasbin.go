package mycasbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	// "github.com/casbin/casbin/v2/log"
	"github.com/casbin/casbin/v2/model"
	"gorm.io/gorm"

	gormAdapter "objgo/team/gorm-adapter/v3"
)

// Initialize the model from a string.
var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

func Setup(db *gorm.DB, _ string) *casbin.SyncedEnforcer {
	Apter, err := gormAdapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	m, err := model.NewModelFromString(text)
	if err != nil {
		panic(err)
	}
	fmt.Println("model.NewMode", m)

	e, err := casbin.NewSyncedEnforcer(m, Apter)
	if err != nil {
		panic(err)
	}
	err = e.LoadPolicy()
	if err != nil {
		panic(err)
	}

	// log.SetLogger()
	e.EnableLog(true)
	return e
}
