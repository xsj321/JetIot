package main

import (
	"JetIot/accountSystem"
	"JetIot/conf"
	"JetIot/model"
	"JetIot/thing"
	"JetIot/util/Log"
	"JetIot/util/mqtt"
	"JetIot/util/mysql"
	"JetIot/util/redis"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
)

func init() {
	conf.InitConfig("conf/defaultConfig.json")
	mysql.InitDB()
	redis.InitRedis()
	mqtt.InitMqttClient()
	mqtt.Publish("server/status", "restart")
}

func runHttpServer() {
	engine := gin.Default()

	//物管理
	thingGroup := engine.Group("thing")
	thingGroup.GET("test", thing.Test)
	thingGroup.POST("register", thing.RegisterThing)
	thingGroup.POST("setThingComponentValue", thing.SetThingComponentValue)
	thingGroup.POST("getThingComponentValue", thing.GetThingComponentValue)

	//账号管理
	accountGroup := engine.Group("account")
	accountGroup.POST("login", accountSystem.Login)
	accountGroup.POST("register", accountSystem.Register)
	err := engine.Run()

	if err != nil {
		Log.E()(err.Error())
	}
}

func runMqttServer() {
	mqtt.Subscribe("thing/function/register", thing.RegisterThingByMqttHandle)
	mqtt.RegisterEventHandle(model.EVENT_COMPONENT_CHANGE_VALUE, "更新设备数据", thing.SetThingComponentValueByMqttHandle)
}

func main() {
	go runHttpServer()
	go runMqttServer()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
