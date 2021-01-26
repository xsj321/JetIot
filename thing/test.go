package thing

import (
	"JetIot/conf"
	"JetIot/util/Log"
	"JetIot/util/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
	//"github.com/jinzhu/gorm"
)

var DbConn *gorm.DB

type Account struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	CreateTime string `json:create_time`
}

func inita() {
	var err error
	mysqlUrl := "root:" + conf.Default.MysqlPassword + "@tcp(" + conf.Default.MysqlServer + ":3306)/cloudDB?charset=utf8mb4&parseTime=true"
	DbConn, err = gorm.Open("mysql", mysqlUrl)
	if err != nil {
		panic(err)
	}
	//DbConn.AutoMigrate(&Post{}, &Comment{})
}

func Test(context *gin.Context) {
	inita()
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

	DbConn.Create(&account)
	fmt.Println(account)
}
