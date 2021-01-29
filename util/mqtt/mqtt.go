package mqtt

import (
	"JetIot/conf"
	"JetIot/util/Log"
	mq "github.com/eclipse/paho.mqtt.golang"
)

var Client mq.Client

func InitMqttClient() {
	opts := mq.NewClientOptions()
	opts.AddBroker("tcp://"+conf.Default.MqttServer+":"+conf.Default.MqttPort).
		SetClientID(conf.Default.MqttClientID).
		SetUsername(conf.Default.MqttUserName).
		SetPassword(conf.Default.MqttPassword).
		SetWill("server/will", "lose_connect", 2, true)

	Client = mq.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		Log.E()("初始化Mqtt客户端失败" + token.Error().Error())
		return
	}
	Log.I()("初始化Mqtt客户端完成")
}

func Publish(topic string, msg string) {
	Log.D()("topic：" + topic + " msg：" + msg)
	token := Client.Publish(topic, 1, false, msg)
	if token.Error() != nil {
		Log.E()("发布消息失败" + token.Error().Error())
		return
	}
	if !token.Wait() {
		Log.E()("发布发送等待")
		return
	}
}

func Subscribe(topic string, callback mq.MessageHandler) {
	if token := Client.Subscribe(topic, 0, callback); token.Wait() && token.Error() != nil {
		Log.E()("订阅失败" + token.Error().Error())
	}
}

func Unsubscribe(topic string) {
	if token := Client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		Log.E()("取消订阅失败" + token.Error().Error())
	}
}
