package router

import (
	"fmt"
	log "objgo/team/core/logger"
	"objgo/team/core/sdk"
	"os"

	"github.com/gin-gonic/gin"
)

func InitRouter() {
	var r *gin.Engine
	h := sdk.Runtime.GetEngine()
	if h == nil {
		log.Fatal("not found engine...")
		os.Exit(-1)
	}
	switch h.(type) {
	case *gin.Engine:
		r = h.(*gin.Engine)
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}
	fmt.Println(r)
	//the jwt middleware
	// authMiddleware, err := common.AuthInit()
	// if err != nil {
	// 	log.Fatalf("JWT Init Error, %s", err.Error())
	// }

	// 注册系统路由

	// 注册业务路由
}
