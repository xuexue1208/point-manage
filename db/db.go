package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/schema"
	"log"
	"time"

	"point-manage/config"
	logs "github.com/wonderivan/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//var (
//	isInit bool
//	GORM   *gorm.DB
//	SqlDB  *sql.DB
//	err    error
//)

//db的初始化函数，与数据库建立连接
//func Init() {
//
//	//判断是否已经初始化了
//	if isInit {
//		return
//	}
//	//组装连接配置
//	//parseTime是查询结果是否自动解析为时间
//	//loc是Mysql的时区设置
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
//		config.Instance.DbInfos.DbUser,
//		config.Instance.DbInfos.DbPwd,
//		config.Instance.DbInfos.DbHost,
//		config.Instance.DbInfos.DbPort,
//		config.Instance.DbInfos.DbName)
//
//	//v2
//	GORM, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: logger.Default.LogMode(logger.Info),
//	})
//	SqlDB, _ = GORM.DB()
//	if err != nil {
//		log.Fatal(err.Error())
//		defer SqlDB.Close()
//
//	} else {
//		fmt.Printf("初始化数据库OK\n")
//		SqlDB.SetMaxIdleConns(10)
//		SqlDB.SetMaxOpenConns(100)
//		SqlDB.SetConnMaxLifetime(time.Hour)
//	}
//
//	isInit = true
//	//GORM.AutoMigrate(
//	//	model.Attribute{},
//	//	model.Classification{},
//	//	model.ClassEvent{},
//	//	model.Demand{},
//	//	model.Event{},
//	//	model.Value{},
//	//	model.Version{},
//	//	model.Role{},
//	//	model.Tags{},
//	//	model.Oplogs{},
//	//	model.User{},
//	//)
//}
var isInit bool
var dbs = make(map[string]*gorm.DB)

func getDB(key string) *gorm.DB {
	if db, ok := dbs[key]; ok {
		return db
	}
	return &gorm.DB{}
}
func GetXadminDB() *gorm.DB {
	return getDB("xadmin")
}
func GetPointDB() *gorm.DB {
	return getDB("point")
}
func NewInit() {
	if isInit {
		return
	}
	for key, dbConfig := range config.Instance.DbInfos {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			dbConfig.DbUser,
			dbConfig.DbPwd,
			dbConfig.DbHost,
			dbConfig.DbPort,
			dbConfig.DbName)
		GORM, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			sqlDB, err := GORM.DB()
			if err == nil {
				sqlDB.SetMaxOpenConns(40)
				sqlDB.SetMaxIdleConns(40)
				sqlDB.SetConnMaxIdleTime(1 * time.Hour)
				logs.Info(string(key), "% 初始化数据库OK\n")
			}
		} else {
			log.Fatal("初始化数据库失败", err.Error())
		}
		//if string(key) == "point" {
		//	GORM.AutoMigrate(
		//		model.Attribute{},
		//		model.Classification{},
		//		model.ClassEvent{},
		//		model.Demand{},
		//		model.Event{},
		//		model.Value{},
		//		model.Version{},
		//		model.Role{},
		//		model.Tags{},
		//		model.Oplogs{},
		//	)
		//}

		dbs[string(key)] = GORM
	}
	isInit = true

}
