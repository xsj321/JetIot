package appServer

import (
	"JetIot/model/account"
	"JetIot/model/response"
	"JetIot/util/Log"
	"JetIot/util/errorCode"
	"JetIot/util/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/**
 * @Description: 登录接口
 * @param context gin上下文
 */
func Login(context *gin.Context) {
	param := account.Account{}
	err := context.ShouldBindJSON(&param)
	if err != nil {
		Log.E()("参数解析错误")
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}
	res, err := mysql.Find("accounts", &param)
	if err != nil {
		Log.E()("数据库查询错误", err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"数据库查询错误或账号密码错误",
			errorCode.ERR_MYSQL_FAILED,
		))
		return
	}
	resTrue := res.(*account.Account)
	if resTrue.Password == param.Password && resTrue.Account == resTrue.Account {
		context.JSON(http.StatusOK, response.GetSuccessResponses("login_succeed"))
		return
	}
	context.JSON(http.StatusOK, response.GetFailResponses(
		"login_error",
		errorCode.ERR_MYSQL_FAILED,
	))
}

func Register(context *gin.Context) {
	param := account.Account{}
	param.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	err := context.ShouldBindJSON(&param)
	if err != nil {
		Log.E()("参数解析错误")
		context.JSON(http.StatusOK, response.GetFailResponses(
			"参数解析错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	res := account.Account{}
	err = mysql.Conn.Table("accounts").Where("account = ?", param.Account).Scan(&res).Error
	if err != nil {
		Log.I()(err)
	}
	if res.Account != "" {
		Log.D()(res.Account)
		context.JSON(http.StatusOK, response.GetFailResponses(
			"注册失败，账号已存在",
			errorCode.ERR_MYSQL_FAILED,
		))
		return
	}
	err = mysql.Create(param)
	if err != nil {
		Log.E()(err.Error())
		context.JSON(http.StatusOK, response.GetFailResponses(
			"注册失败，系统错误",
			errorCode.ERR_SVR_INTERNAL,
		))
		return
	}

	context.JSON(http.StatusOK, response.GetSuccessResponses("reg_succeed"))

}
