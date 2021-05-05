package thingServer

import (
	"JetIot/model/event"
	"JetIot/model/response"
	"JetIot/model/thingModel"
	"JetIot/util"
	"JetIot/util/Log"
	"JetIot/util/errorCode"
	"JetIot/util/mqtt"
	"JetIot/util/redis"
	"encoding/json"
	mq "github.com/eclipse/paho.mqtt.golang"
)

/**
 * @Description: 注册回调
 * @param client
 * @param message
 */
func RegisterThingByMqttHandle(client mq.Client, message mq.Message) {
	thing := thingModel.Thing{}
	err := mqtt.ShouldBindJSON(message, &thing)
	replayTopic := "thingServer/function/register/" + thing.Id

	if err != nil {
		Log.E()("参数解析错误")
		mqtt.Replay(replayTopic, response.GetMqttFailResponses(
			"参数解析错误",
			event.EVENT_OTHER,
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	marshal, _ := json.Marshal(thing)

	//保存到缓存库
	err = redis.Set("thingServer:"+thing.Id, string(marshal))
	if err != nil {
		Log.E()("保存到缓存库错误" + err.Error())
		mqtt.Replay(replayTopic, response.GetMqttFailResponses(
			"保存到缓存库错误",
			event.EVENT_OTHER,
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	mqtt.Replay(replayTopic, response.GetMqttSuccessResponses("注册完成", event.EVENT_OTHER))
}

/**
 * @Description: 设置解析并存储来自设备的值
 * @param context
 */
func SetThingComponentValueByMqttHandle(client mq.Client, message mq.Message) {
	ov := thingModel.ThingComponentValueOV{}
	err := mqtt.ShouldBindJSON(message, &ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		return
	}

	thing, err := util.LoadThing(ov.Id)
	if err != nil {
		Log.E()("加载错误" + err.Error())
		return
	}

	thing.Components[ov.ComponentName].Do(ov.Value)
	err = util.Commit(thing)
	if err != nil {
		Log.E()("更新缓存错误" + err.Error())
		return
	}
	if hook, has := ComponentHooks[ov.ComponentName]; has {
		Log.D()("执行Hook:", hook.ComponentName)
		hook.ComponentHookFuc(ov)
	} else {
		Log.I()("Hook:", hook.ComponentName, "ID:", ov.Id, "未找到")
	}

}

/**
 * @Description: 设备上线回调
 * @param client mqtt客户端
 * @param message 消息对象
 */
func DeviceOnlineByMqttHandle(client mq.Client, message mq.Message) {
	ov := thingModel.ThingInitOV{}
	err := mqtt.ShouldBindJSON(message, &ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		return
	}

	isReg := util.IsRegisterThing(ov.Id)
	if isReg {
		Log.D()("设备已经注册")
		err := util.SetDeviceOnlineStatus(ov.Id, true)
		if err != nil {
			mqtt.ReplayToDevice(ov.Id, response.GetMqttFailResponses("上线失败", ov.EventId, errorCode.ERR_SVR_INTERNAL))
		}
		mqtt.ReplayToDevice(ov.Id, response.GetMqttSuccessResponses("设备已经注册，并上线", ov.EventId, nil))

		return
	} else {
		Log.D()("设备未注册")
		mqtt.ReplayToDevice(ov.Id, response.GetMqttFailResponses("设备未注册", ov.EventId, errorCode.ERR_DEVICE_NOT_FIND))

		return
	}
}

/**
 * @Description: 设备下线回调
 * @param client mqtt客户端
 * @param message 消息对象
 */
func DeviceOfflineByMqttHandle(client mq.Client, message mq.Message) {
	Log.D()("遗嘱内容", string(message.Payload()))
	ov := thingModel.ThingInitOV{}
	err := mqtt.ShouldBindJSON(message, &ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		return
	}

	if util.IsRegisterThing(ov.Id) {
		err := util.SetDeviceOnlineStatus(ov.Id, false)
		if err != nil {
			Log.E()("设备下线失败，但是已经离线", err.Error(), "ID:", ov.Id)
			return
		}
	} else {
		Log.E()("未注册设备")
		return
	}

}

// 设备心跳
func DeviceHeatByMqttHandle(client mq.Client, message mq.Message) {
	Log.D()("心跳内容", string(message.Payload()))
	ov := thingModel.ThingInitOV{}
	err := mqtt.ShouldBindJSON(message, &ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		return
	}

	if util.IsRegisterThing(ov.Id) {
		err := util.UpdateDeviceOnLineStatus(ov.Id)
		if err != nil {
			Log.E()("设备心跳失败", err.Error(), "ID:", ov.Id)
			return
		}
	} else {
		Log.E()("未注册设备")
		return
	}

}
