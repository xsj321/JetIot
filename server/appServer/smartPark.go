package appServer

import (
	"JetIot/model/park"
	"JetIot/model/response"
	"JetIot/util"
	"JetIot/util/Log"
	"JetIot/util/errorCode"
	"JetIot/util/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetCoverListByLocation(context *gin.Context) {
	ov := park.CoverOV{}
	err := context.ShouldBindJSON(&ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	Log.D()("请求位置为:", ov.Place)

	if ov.Place == "" || ov.Place == "null" {
		Log.D()("未使用地址搜索")
		ov.Place = "%"
	}

	coverList := make([]park.Cover, 0)
	err = mysql.Conn.Table("cover_status").Where("place like ?", ov.Place).Scan(&coverList).Error
	if err != nil {
		Log.E()("数据库查询错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"数据库查询错误",
			errorCode.ERR_MYSQL_FAILED,
		))
		return
	}

	waring := make([]park.Cover, 0)
	detail := make([]park.Cover, 0)

	for _, cover := range coverList {
		detail = append(detail, cover)
		if cover.Waring == 1 {
			waring = append(waring, cover)
		}
	}

	resOV := park.CoverResOV{
		WaringList: waring,
		Detail:     detail,
	}

	context.JSON(
		http.StatusOK,
		response.GetSuccessResponses("获取数据列表成功", resOV),
	)

}

/**
 * @Description: 获取当前井盖的状态
 * @param id 设备ID
 * @return bool 井盖状态 true：正常 false：异常
 */
func GetCoverStatus(id string) bool {
	thing, err := util.LoadThing(id)
	if err != nil {
		Log.E()("加载物模型错误" + err.Error())
		return false
	}
	cover := thing.Components["cover"]
	coverStatus := cover.Value.(bool)

	//var coverStatus int
	//err := mysql.Conn.Table("cover_status").Select("waring").Where("device_id = ?", id).Scan(&coverStatus).Error
	//if err != nil {
	//	Log.E()("数据库查询错误" + err.Error())
	//	return false
	//}
	if coverStatus {
		return true
	}
	return false
}

//修复井盖
func FixCoverWaring(context *gin.Context) {
	ov := park.CoverFixOV{}
	err := context.ShouldBindJSON(&ov)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	Log.D()("请求修复ID为:", ov.Id)

	err = mysql.Conn.Table("cover_status").Where("id = ?", ov.Id).Update("waring", 0).Error
	if err != nil {
		Log.E()("数据库更新错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"数据库更新错误",
			errorCode.ERR_MYSQL_FAILED,
		))
		return
	}

	context.JSON(
		http.StatusOK,
		response.GetSuccessResponses("修复井盖成功"),
	)

}
