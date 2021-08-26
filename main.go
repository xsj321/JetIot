package main

import (
	"JetIot/conf"
	"JetIot/server/appServer"
	"JetIot/util/Log"
	"JetIot/util/mysql"
	"JetIot/util/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"os"
	"os/signal"
)

var configPath = pflag.StringP("config", "c", "conf/defaultConfig.json", "the config json file path")

func initEvn() {
	conf.InitConfig(*configPath)
	mysql.InitDB()
	redis.InitRedis()
}

func runHttpServer() {
	engine := gin.Default()

	//app
	appGroup := engine.Group("app")

	thingGroup := appGroup.Group("note")
	thingGroup.POST("saveNote", appServer.SaveNote)
	thingGroup.POST("getNote", appServer.GetNote)
	thingGroup.POST("delNote", appServer.DelNote)
	//账号管理
	accountGroup := engine.Group("account")
	accountGroup.POST("login", appServer.Login)
	accountGroup.POST("register", appServer.Register)
	accountGroup.POST("listAccount", appServer.ListAccount)
	accountGroup.POST("addFriend", appServer.AddFriend)
	accountGroup.POST("acceptFriend", appServer.AcceptFriend)
	accountGroup.POST("getFriendList", appServer.GetFriendList)
	accountGroup.POST("getFriendRequestList", appServer.GetFriendRequestList)
	err := engine.Run(":8888")

	if err != nil {
		Log.E()(err.Error())
	}
}

func main() {
	initEvn()
	go runHttpServer()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
