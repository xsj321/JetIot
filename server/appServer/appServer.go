/**
 * @Description: 通过HTTP对设备进行注册修改等操作
 */
package appServer

import (
	"JetIot/model/note"
	"JetIot/model/response"
	"JetIot/util/Log"
	"JetIot/util/errorCode"
	"JetIot/util/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNote(context *gin.Context) {
	note := note.UserNote{}
	err := context.ShouldBindJSON(&note)
	if err != nil && note.Id == 0 {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	Log.D()("获取的id为：", note.Id)
	err = mysql.Conn.Table("user_note").Where("id = ?", note.Id).Scan(&note).Error
	if err != nil {
		Log.E()("笔记获取错误" + err.Error())
		if err.Error() == "record not found" {
			context.JSON(http.StatusOK, response.GetFailResponses(
				"笔记不存在",
				errorCode.ERR_SVR_INTERNAL,
			))
			return
		}
		context.JSON(http.StatusOK, response.GetFailResponses(
			"笔记获取错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	Log.D()("获取的内容为：", note.NoteData)
	context.JSON(http.StatusOK, response.GetSuccessResponses(
		"笔记获取成功",
		note,
	))
}

func DelNote(context *gin.Context) {
	note := note.UserNote{}
	err := context.ShouldBindJSON(&note)
	if err != nil && note.Id == 0 {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	Log.D()("删除的id为：", note.Id)
	err = mysql.Conn.Table("user_note").Delete(&note).Error
	if err != nil {
		Log.E()("笔记删除错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"笔记删除错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	context.JSON(http.StatusOK, response.GetSuccessResponses(
		"笔记删除成功",
	))
}

func SaveNote(context *gin.Context) {
	note := note.UserNote{}
	err := context.ShouldBindJSON(&note)
	if err != nil {
		Log.E()("参数解析错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	err = mysql.Conn.Table("user_note").Create(&note).Error
	if err != nil {
		Log.E()("笔记保存错误" + err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"笔记保存错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	Log.D()("插入的id为：", note.Id)
	context.JSON(http.StatusOK, response.GetSuccessResponses(
		"笔记保存成功",
		struct {
			NoteId int `json:"note_id"`
		}{
			NoteId: note.Id,
		},
	))
}
