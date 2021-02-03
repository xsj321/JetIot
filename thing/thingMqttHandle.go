package thing

import (
	"JetIot/model"
	"JetIot/model/event"
	"JetIot/model/response"
	"JetIot/util/Log"
	"JetIot/util/errorCode"
	"JetIot/util/mqtt"
	"JetIot/util/redis"
	"encoding/json"
	mq "github.com/eclipse/paho.mqtt.golang"
)

func RegisterThingByMqttHandle(client mq.Client, message mq.Message) {
	thing := model.Thing{}
	err := mqtt.ShouldBindJSON(message, &thing)
	replayTopic := "thing/function/register/" + thing.Id

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
	err = redis.Set("thing:"+thing.Id, string(marshal))
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
	ov := model.ThingComponentValueOV{}
	err := mqtt.ShouldBindJSON(message, &ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		return
	}

	thing, err := LoadThing(ov.Id)
	if err != nil {
		Log.E()("加载错误" + err.Error())
		return
	}

	thing.Components[ov.ComponentName].Do(ov.Value)
	err = commit(thing)
	if err != nil {
		Log.E()("更新缓存错误" + err.Error())
		return
	}
}

/**
 * @Description: 设备上线回调
 * @param client mqtt客户端
 * @param message 消息对象
 */
func DeviceOnlineByMqttHandle(client mq.Client, message mq.Message) {
	ov := model.ThingInitOV{}
	err := mqtt.ShouldBindJSON(message, &ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		return
	}

	isReg := isRegisterThing(ov.Id)
	if isReg {
		Log.D()("设备已经注册")
		err := setDeviceOnlineStatus(ov.Id, true)
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
	ov := model.ThingInitOV{}
	err := mqtt.ShouldBindJSON(message, &ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		return
	}

	if isRegisterThing(ov.Id) {
		err := setDeviceOnlineStatus(ov.Id, false)
		if err != nil {
			Log.E()("设备下线失败，但是已经离线", err.Error(), "ID:", ov.Id)
			return
		}
	} else {
		Log.E()("未注册设备")
		return
	}

}
