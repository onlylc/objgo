package database

import (
	"time"

	log "objgo/team/core/logger"
	"objgo/team/core/sdk"

	toolsConfig "objgo/team/core/sdk/config"
	"objgo/team/core/sdk/pkg"
	mycasbin "objgo/team/core/sdk/pkg/casbin"
	toolsDB "objgo/team/core/tools/database"
	. "objgo/team/core/tools/gorm/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"objgo/common/global"
)

// Setup 配置数据
func Setup() {
	for k := range toolsConfig.DatabasesConfig {
		setupSimpleDatabase(k, toolsConfig.DatabasesConfig[k])
	}
}

func setupSimpleDatabase(host string, c *toolsConfig.Database) {
	if global.Driver == "" {
		global.Driver = c.Driver
	}
	log.Infof("%s => %s", host, pkg.Green(c.Source))
	registers := make([]toolsDB.ResolverConfigure, len(c.Registers))
	for i := range c.Registers {
		registers[i] = toolsDB.NewResolverConfigure(
			c.Registers[i].Sources,
			c.Registers[i].Replicas,
			c.Registers[i].Policy,
			c.Registers[i].Tables)

	}

	resolverConfig := toolsDB.NewConfigure(c.Source, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxIdleTime, c.ConnMaxLifeTime, registers)
	db, err := resolverConfig.Init(&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: New(
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					log.DefaultLogger.Options().Level.LevelForGorm()),
			},
		),
	}, opens[c.Driver])

	if err != nil {
		log.Fatal(pkg.Red(c.Driver+" connect error :"), err)
	} else {
		log.Info(pkg.Green(c.Driver + " connect success !"))
	}

	e := mycasbin.Setup(db, "sys_")

	sdk.Runtime.SetDb(host, db)
	sdk.Runtime.SetCasbin(host, e)
}
