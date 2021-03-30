/**
 * @Description: 通过HTTP对设备进行注册修改等操作
 */
package appServer

import (
	"JetIot/model/account"
	"JetIot/model/response"
	"JetIot/model/thingModel"
	"JetIot/util"
	"JetIot/util/Log"
	"JetIot/util/errorCode"
	"JetIot/util/mysql"
	"JetIot/util/redis"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * @Description: 注册物类型与组件
 * @param context
 */
func RegisterThing(context *gin.Context) {
	thing := thingModel.Thing{}
	err := context.ShouldBindJSON(&thing)
	if err != nil {
		Log.E()("参数解析错误")
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	marshal, _ := json.Marshal(thing)

	//保存到缓存库
	err = redis.Set("thingServer:"+thing.Id, string(marshal))
	if err != nil {
		Log.E()("保存到缓存库错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"保存到缓存库错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	context.JSON(http.StatusOK, response.GetSuccessResponses("注册完成"))
}

/**
 * @Description: 设置组件的值
 * @param context
 */
func SetThingComponentValue(context *gin.Context) {
	ov := thingModel.ThingComponentValueOV{}
	err := context.ShouldBindJSON(&ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	thing, err := util.LoadThing(ov.Id)
	if err != nil {
		Log.E()("加载错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"加载错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	thing.Components[ov.ComponentName].Do(ov.Value)
	err = util.Commit(thing)
	if err != nil {
		Log.E()("更新缓存错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"更新缓存错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	context.JSON(http.StatusOK, response.GetSuccessResponses("修改完成"))
}

/**
 * @Description: 解绑设备
 * @param context
 */
func UnBindingDevice(context *gin.Context) {
	deviceOV := account.BindingDevice{}
	err := context.ShouldBindJSON(&deviceOV)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	err = redis.Del("thing:" + deviceOV.DeviceId)
	if err != nil {
		Log.E()("删除设备失败", err)
		context.JSON(http.StatusOK, response.GetFailResponses(
			"删除设备失败",
			errorCode.ERR_REDIS_FAILED,
		))
		return
	}

	err = mysql.Conn.Delete(&deviceOV).Error
	if err != nil {
		Log.E()("删除设备失败", err)
		context.JSON(http.StatusOK, response.GetFailResponses(
			"删除设备失败",
			errorCode.ERR_MYSQL_FAILED,
		))
		return
	}

	context.JSON(http.StatusOK, response.GetSuccessResponses(
		"设备解绑完成",
		errorCode.ERR_NO_ERROR,
	))
}

/**
 * @Description: 绑定设备
 * @param context
 */
func BindingDevice(context *gin.Context) {
	deviceOV := account.BindingDevice{}
	err := context.ShouldBindJSON(&deviceOV)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	//验证设备是否被绑定
	isBinding := util.IsDeviceBinding(deviceOV.DeviceId)

	if isBinding {
		context.JSON(http.StatusOK, response.GetFailResponses(
			"设备已经被绑定",
			errorCode.ERR_DEV_BINDED,
		))

		return
	}

	deviceOV.CreateTime = mysql.GetFormatNowTime()
	err = mysql.Conn.Create(&deviceOV).Error
	if err != nil {
		Log.E()("设备保存错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"设备保存错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	context.JSON(http.StatusOK, response.GetSuccessResponses(
		"设备绑定完成",
		errorCode.ERR_NO_ERROR,
	))

}

func GetDeviceInfoListByAccount(context *gin.Context) {
	ov := account.Account{}
	err := context.ShouldBindJSON(&ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	var deviceList []account.BindingDevice
	err = mysql.Conn.Select("device_id").Table("binding_devices").Where("account = ?", ov.Account).Scan(&deviceList).Error
	if err != nil {
		Log.E()("查询设备列表失败", err)
		context.JSON(http.StatusOK, response.GetFailResponses(
			"查询设备列表失败",
			errorCode.ERR_MYSQL_FAILED,
		))
		return
	}
	if len(deviceList) == 0 {
		Log.D()("设备列表为空", err)
		context.JSON(http.StatusOK, response.GetFailResponses(
			"设备列表为空",
			errorCode.ERR_DEVICE_LIST_EMPTY,
		))
		return
	}

	var thingList []thingModel.Thing

	for _, s := range deviceList {
		thing, err := util.LoadThing(s.DeviceId)
		if err != nil {
			Log.E()("加载错误" + err.Error())
			context.JSON(http.StatusOK, response.GetFailResponses(
				"加载错误",
				errorCode.ERR_SVR_INTERNAL,
			))
			return
		}
		thingList = append(thingList, *thing)
	}

	context.JSON(http.StatusOK, response.GetSuccessResponses(
		"设备列表获取成功",
		thingList,
	))
}

func GetThingComponentValue(context *gin.Context) {
	ov := thingModel.ThingComponentValueOV{}
	err := context.ShouldBindJSON(&ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	thing, err := util.LoadThing(ov.Id)
	if err != nil {
		Log.E()("加载错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"加载错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	res := thing.Components[ov.ComponentName].Call()
	ov.Value = res
	context.JSON(http.StatusOK, response.GetSuccessResponses("查询成功", ov))
}

func GetThingAllValue(context *gin.Context) {
	ov := thingModel.ThingInitOV{}
	err := context.ShouldBindJSON(&ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	thing, err := util.LoadThing(ov.Id)
	if err != nil {
		Log.E()("加载错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"加载错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	res, _ := json.Marshal(thing)
	//ov.Value = res
	context.JSON(http.StatusOK, response.GetSuccessResponses("查询成功", string(res)))
}
