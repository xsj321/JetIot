package main

import (
	"JetIot/conf"
	"JetIot/model/event"
	"JetIot/server/appServer"
	"JetIot/server/thingServer"
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
	thingServer.InitComponentHooks()
	mqtt.Publish("server/status", "restart")
}

func runHttpServer() {
	engine := gin.Default()

	//app
	appGroup := engine.Group("app")

	//物管理
	thingGroup := appGroup.Group("thing")
	thingGroup.POST("register", appServer.RegisterThing)
	thingGroup.POST("setThingComponentValue", appServer.SetThingComponentValue)
	thingGroup.POST("getThingComponentValue", appServer.GetThingComponentValue)
	thingGroup.POST("getThingAllValue", appServer.GetThingAllValue)
	thingGroup.POST("bindingDevice", appServer.BindingDevice)
	thingGroup.POST("unbindingDevice", appServer.UnBindingDevice)
	thingGroup.POST("getDeviceListByAccount", appServer.GetDeviceInfoListByAccount)

	//智慧园区专属API
	thingGroup.POST("getCoverListByLocation", appServer.GetCoverListByLocation)
	thingGroup.POST("fixCoverWaring", appServer.FixCoverWaring)

	//账号管理
	accountGroup := engine.Group("account")
	accountGroup.POST("login", appServer.Login)
	accountGroup.POST("register", appServer.Register)
	accountGroup.POST("listAccount", appServer.ListAccount)
	accountGroup.POST("addFriend", appServer.AddFriend)
	accountGroup.POST("acceptFriend", appServer.AcceptFriend)
	accountGroup.POST("getFriendList", appServer.GetFriendList)
	accountGroup.POST("getFriendRequestList", appServer.GetFriendRequestList)
	err := engine.Run()

	if err != nil {
		Log.E()(err.Error())
	}
}

func runMqttServer() {
	mqtt.Subscribe("thingServer/function/register", thingServer.RegisterThingByMqttHandle)
	mqtt.RegisterEventHandle(event.EVENT_COMPONENT_CHANGE_VALUE, "更新设备组件数据", thingServer.SetThingComponentValueByMqttHandle)
	mqtt.RegisterEventHandle(event.EVENT_THING_DEVICE_ONLIONE, "设备上线初始化", thingServer.DeviceOnlineByMqttHandle)
	mqtt.RegisterEventHandle(event.EVENT_THING_DEVICE_OFFLIONE, "设备离线遗嘱", thingServer.DeviceOfflineByMqttHandle)
	mqtt.RegisterEventHandle(event.EVENT_THING_DEVICE_HEATBEAT, "设备心跳", thingServer.DeviceHeatByMqttHandle)
	mqtt.EventListenStart()

	//添加事件钩子

	//井盖事件钩子
	thingServer.AddComponentHook("cover", thingServer.UpdateCoverStatusHook)
}

func main() {
	initEvn()
	go runHttpServer()
	go runMqttServer()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
