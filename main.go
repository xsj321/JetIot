package main

import (
	"JetIot/accountSystem"
	"JetIot/conf"
	"JetIot/thing"
	"JetIot/util/mysql"
	"JetIot/util/redis"
	"github.com/gin-gonic/gin"
)

func init() {
	conf.InitConfig("conf/defaultConfig.json")
	mysql.InitDB()
	redis.InitRedis()
}

func main() {
	engine := gin.Default()

	thingGroup := engine.Group("thing")
	thingGroup.GET("test", thing.Test)
	thingGroup.POST("register", thing.RegisterThing)
	thingGroup.POST("setThingComponentValue", thing.SetThingComponentValue)
	thingGroup.POST("getThingComponentValue", thing.GetThingComponentValue)

	engine.POST("/login", accountSystem.Login)
	engine.POST("/register", accountSystem.Register)
	err := engine.Run()
	if err != nil {
		panic(err)
	}
}
