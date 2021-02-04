/**
 * @Description: 通过HTTP对设备进行注册修改等操作
 */
package thing

import (
	"JetIot/model/account"
	"JetIot/model/response"
	"JetIot/model/thingModel"
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
	err = redis.Set("thing:"+thing.Id, string(marshal))
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

	thing, err := LoadThing(ov.Id)
	if err != nil {
		Log.E()("加载错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"加载错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	thing.Components[ov.ComponentName].Do(ov.Value)
	err = commit(thing)
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
	isBinding := isDeviceBinding(deviceOV.DeviceId)

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

	thing, err := LoadThing(ov.Id)
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
	context.JSON(http.StatusOK, response.GetSuccessResponses("修改完成", ov))
}
