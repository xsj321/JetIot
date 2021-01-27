package mysql

import (
	"JetIot/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Conn *gorm.DB

func InitDB() {
	var err error
	mysqlUrl := "root:" + conf.Default.MysqlPassword + "@tcp(" + conf.Default.MysqlServer + ":3306)/cloudDB?charset=utf8mb4&parseTime=true"
	Conn, err = gorm.Open("mysql", mysqlUrl)
	if err != nil {
		panic(err)
	}
	//DbConn.AutoMigrate(&Post{}, &Comment{})
}

func Create(model interface{}) error {
	create := Conn.Create(model)
	err := create.Error
	return err
}

/**
 * @Description: 根据条件查询数据库
 * @param table	 表名
 * @param query	 查询用的对象
 * @param columns 筛选的字段，可以为空则查找全部
 * @return interface{} 返回的空接口中为对象的指针
 * @return error
 */
func Find(table string, query interface{}, columns ...string) (interface{}, error) {
	var err error
	if len(columns) == 0 {
		err = Conn.Table(table).Where(query).Scan(query).Error
	} else {
		err = Conn.Table(table).Select(columns).Where(query).Scan(query).Error
	}
	return query, err
}
