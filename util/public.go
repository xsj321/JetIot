package util

import (
	"JetIot/model/account"
	"JetIot/model/thingModel"
	"JetIot/util/Log"
	"JetIot/util/mysql"
	"JetIot/util/redis"
	"encoding/json"
)

func LoadThing(id string) (*thingModel.Thing, error) {
	thing := thingModel.Thing{}
	get, err := redis.Get("thingServer:" + id)
	if err != nil {
		Log.E()("查找物体错误" + err.Error())
		return &thing, err
	}

	err = json.Unmarshal(get, &thing)
	if err != nil {
		Log.E()("解析物体错误" + err.Error())
		return &thing, err
	}

	return &thing, nil
}

func Commit(thing *thingModel.Thing) error {
	marshal, _ := json.Marshal(*thing)
	Log.D()(string(marshal))
	//保存到缓存库
	err := redis.Set("thingServer:"+thing.Id, string(marshal))
	if err != nil {
		Log.E()("保存到缓存库错误" + err.Error())
		return err
	}
	return nil
}

func IsRegisterThing(id string) bool {
	_, err := redis.Get("thingServer:" + id)
	if err != nil {
		return false
	}
	return true
}

/**
 * @Description: 设置设备在线状态 在线：1 离线：0
 * @param id 设备id
 * @param isOnline 是否在线
 */
func SetDeviceOnlineStatus(id string, isOnline bool) error {
	if isOnline {
		err := redis.Set("thingServer:online_status:"+id, 1)
		if err != nil {
			Log.E()(id, "更新设备状态错误", err)
			return err
		}
	} else {
		err := redis.Set("thingServer:online_status:"+id, 0)
		if err != nil {
			Log.E()(id, "更新设备状态错误", err)
			return err
		}
	}
	return nil
}

/**
 * @Description: 查询设备是否绑定
 * @param id 设备ID
 * @return bool 是否被绑定
 */
func IsDeviceBinding(id string) bool {
	device := account.BindingDevice{}
	err := mysql.Conn.Select("account").Table("binding_devices").Where("device_id = ?", id).Scan(&device).Error
	if err != nil {
		Log.I()("未找到数据，设备", id, "未绑定", err)
		return false
	} else {
		Log.I()("设备", id, "已经绑定用户：", device.Account)
		return true
	}
}
