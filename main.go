package main

import (
	"JetIot/conf"
	"JetIot/thing"
	"github.com/gin-gonic/gin"
)

func init() {
	conf.InitConfig("conf/defaultConfig.json")
}

func main() {

	engine := gin.Default()
	engine.GET("/thing/test", thing.Test)
	engine.Run()
}
