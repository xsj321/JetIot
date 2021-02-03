/**
 * @Description: 通过HTTP对设备进行注册修改等操作
 */
package thing

import (
	"JetIot/model"
	"JetIot/model/response"
	"JetIot/util/Log"
	"JetIot/util/errorCode"
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
	thing := model.Thing{}
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
	ov := model.ThingComponentValueOV{}
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

func GetThingComponentValue(context *gin.Context) {
	ov := model.ThingComponentValueOV{}
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
