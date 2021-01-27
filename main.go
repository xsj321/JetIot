package main

import (
	"JetIot/conf"
	"JetIot/thing"
	"JetIot/util/mysql"
	"github.com/gin-gonic/gin"
)

func init() {
	conf.InitConfig("conf/defaultConfig.json")
}

func main() {
	mysql.InitDB()
	engine := gin.Default()
	engine.GET("/thing/test", thing.Test)
	engine.Run()
}
