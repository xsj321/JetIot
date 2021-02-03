package main

import (
	"JetIot/accountSystem"
	"JetIot/conf"
	"JetIot/model/event"
	"JetIot/thing"
	"JetIot/util/Log"
	"JetIot/util/mqtt"
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
	mqtt.RegisterEventHandle(event.EVENT_COMPONENT_CHANGE_VALUE, "更新设备数据", thing.SetThingComponentValueByMqttHandle)
	mqtt.RegisterEventHandle(event.EVENT_THING_DEVICE_ONLIONE, "设备上线初始化", thing.DeviceOnlineByMqttHandle)
	mqtt.RegisterEventHandle(event.EVENT_THING_DEVICE_OFFLIONE, "设备离线遗嘱", thing.DeviceOfflineByMqttHandle)
	mqtt.EventListenStart()
}

func main() {
	initEvn()
	go runHttpServer()
	go runMqttServer()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
