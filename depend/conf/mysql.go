package conf

import (
	"time"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/go-sql-driver/mysql"
)

// mysql
func NewMysqlPool(poolcfg *MysqlConf) *gorm.DB {

	dsn := mysql.Config{
		Addr:    poolcfg.ServerAddress,
		User:    poolcfg.UserName,
		Passwd:  poolcfg.Password,
		Net:     "tcp",
		DBName:  poolcfg.DBName,
		Params:  map[string]string{"charset": "utf8", "parseTime": "True", "loc": "Local"},
		Timeout: time.Duration(5 * time.Second),
	}

	db, err := gorm.Open("mysql", dsn.FormatDSN())
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	// SetMaxIdleConns 设置空闲连接池中的最大连接数。
	db.DB().SetMaxIdleConns(poolcfg.MaxIdleConn)

	// SetMaxOpenConns 设置数据库连接最大打开数。
	db.DB().SetMaxOpenConns(poolcfg.MaxOpenConn)

	// SetConnMaxLifetime 设置可重用连接的最长时间
	//db.DB().SetConnMaxLifetime(time.Hour)
	return db
}