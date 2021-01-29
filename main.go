package main

import (
	"JetIot/accountSystem"
	"JetIot/conf"
	"JetIot/thing"
	"JetIot/util/Log"
	"JetIot/util/mqtt"
	"JetIot/util/mysql"
	"JetIot/util/redis"
	"github.com/gin-gonic/gin"
)

func init() {
	conf.InitConfig("conf/defaultConfig.json")
	mysql.InitDB()
	redis.InitRedis()
	mqtt.InitMqttClient()
	mqtt.Publish("server/status", "restart")
}

func main() {
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
