/**
 * @Description: 通过HTTP对设备进行注册修改等操作
 */
package thing

import (
	"JetIot/model"
	"JetIot/util"
	"JetIot/util/Log"
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
		context.JSON(http.StatusOK, model.GetFailResponses(
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
		context.JSON(http.StatusOK, model.GetFailResponses(
			"保存到缓存库错误",
			util.ERR_SVR_INTERNAL,
		))
		return
	}
	context.JSON(http.StatusOK, model.GetSuccessResponses("注册完成"))
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
		context.JSON(http.StatusOK, model.GetFailResponses(
			"参数解析错误",
			util.ERR_SVR_INTERNAL,
		))
		return
	}

	thing, err := LoadThing(ov.Id)
	if err != nil {
		Log.E()("加载错误" + err.Error())
		context.JSON(http.StatusOK, model.GetFailResponses(
			"加载错误",
			util.ERR_SVR_INTERNAL,
		))
		return
	}

	thing.Components[ov.ComponentName].Do(ov.Value)
	err = commit(thing)
	if err != nil {
		Log.E()("更新缓存错误" + err.Error())
		context.JSON(http.StatusOK, model.GetFailResponses(
			"更新缓存错误",
			util.ERR_SVR_INTERNAL,
		))
		return
	}
	context.JSON(http.StatusOK, model.GetSuccessResponses("修改完成"))
}

func GetThingComponentValue(context *gin.Context) {
	ov := model.ThingComponentValueOV{}
	err := context.ShouldBindJSON(&ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, model.GetFailResponses(
			"参数解析错误",
			util.ERR_SVR_INTERNAL,
		))
		return
	}

	thing, err := LoadThing(ov.Id)
	if err != nil {
		Log.E()("加载错误" + err.Error())
		context.JSON(http.StatusOK, model.GetFailResponses(
			"加载错误",
			util.ERR_SVR_INTERNAL,
		))
		return
	}
	res := thing.Components[ov.ComponentName].Call()
	ov.Value = res
	context.JSON(http.StatusOK, model.GetSuccessResponses("修改完成", ov))
}

func LoadThing(id string) (*model.Thing, error) {
	thing := model.Thing{}
	get, err := redis.Get("thing:" + id)
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

func commit(thing *model.Thing) error {
	marshal, _ := json.Marshal(*thing)
	Log.D()(string(marshal))
	//保存到缓存库
	err := redis.Set("thing:"+thing.Id, string(marshal))
	if err != nil {
		Log.E()("保存到缓存库错误" + err.Error())
		return err
	}
	return nil
}
