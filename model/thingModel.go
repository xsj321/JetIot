package model

import "JetIot/util/Log"

type BaseInfo struct {
	EventId int `json:"event_id"`
}

type Thing struct {
	BaseInfo
	Id         string                `json:"id"`
	Name       string                `json:"name"`
	Components map[string]*Component `json:"components"`
}

type Component struct {
	BaseInfo
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type ThingComponentValueOV struct {
	BaseInfo
	Id            string      `json:"id"`             //物体id
	ComponentName string      `json:"component_name"` //组件名称
	Value         interface{} `json:"value"`          //值
}

type ThingInitOV struct {
	BaseInfo
	Id string `json:"id"` //物体id
}

func (c *Component) Do(value interface{}) {
	Log.D()(value)
	c.Value = value
}

func (c *Component) Call() interface{} {
	return c.Value
}
