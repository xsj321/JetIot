package thing

import (
	"JetIot/util/Log"
	"JetIot/util/mysql"
	"JetIot/util/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	//"github.com/jinzhu/gorm"
)

type Account struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	CreateTime string `json:create_time`
}

func Test(context *gin.Context) {
	set := redis.Set("time2", "gg")
	if set != nil {
		Log.D()(set)
	}
	get, _ := redis.Get("time2")
	Log.I()(get)

	format := time.Now().Format("2006-01-02 15:04:05")
	account := Account{
		Account:    "xsjcool@outook.com",
		Name:       "xsjcool",
		Password:   "123",
		Type:       0,
		CreateTime: format,
	}

	mysql.Create(&account)
	fmt.Println(account)
}
