package thing

import (
	"JetIot/model"
	"JetIot/util"
	"JetIot/util/Log"
	"JetIot/util/mqtt"
	"JetIot/util/redis"
	"encoding/json"
	mq "github.com/eclipse/paho.mqtt.golang"
)

func getToDevReplayTopic(devID string) string {
	return "thing/entity/" + devID + "/todevice"
}

func RegisterThingByMqttHandle(client mq.Client, message mq.Message) {
	thing := model.Thing{}
	err := mqtt.ShouldBindJSON(message, &thing)
	replayTopic := "thing/function/register/" + thing.Id

	if err != nil {
		Log.E()("参数解析错误")
		mqtt.Replay(replayTopic, model.GetFailResponses(
			"参数解析错误",
			util.ERR_SVR_INTERNAL,
		))
		return
	}

	marshal, _ := json.Marshal(thing)

	//保存到缓存库
	err = redis.Set("thing:"+thing.Id, string(marshal))
	if err != nil {
		Log.E()("保存到缓存库错误" + err.Error())
		mqtt.Replay(replayTopic, model.GetFailResponses(
			"保存到缓存库错误",
			util.ERR_SVR_INTERNAL,
		))
		return
	}
	mqtt.Replay(replayTopic, model.GetSuccessResponses("注册完成"))
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
