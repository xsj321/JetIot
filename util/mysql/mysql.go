package mysql

import (
	"JetIot/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var dbConn *gorm.DB

func InitDB() {
	var err error
	mysqlUrl := "root:" + conf.Default.MysqlPassword + "@tcp(" + conf.Default.MysqlServer + ":3306)/cloudDB?charset=utf8mb4&parseTime=true"
	dbConn, err = gorm.Open("mysql", mysqlUrl)
	if err != nil {
		panic(err)
	}
	//DbConn.AutoMigrate(&Post{}, &Comment{})
}

func Create(model interface{}) error {
	create := dbConn.Create(model)
	err := create.Error
	return err
}
