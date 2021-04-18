package appServer

import (
	"JetIot/model/park"
	"JetIot/model/response"
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
		if cover.Waring == 0 {
			detail = append(detail, cover)
		} else {
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
