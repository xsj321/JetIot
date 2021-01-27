package main

import (
	"JetIot/accountSystem"
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
	engine.POST("/login", accountSystem.Login)
	engine.POST("/register", accountSystem.Register)
	engine.Run()
}
