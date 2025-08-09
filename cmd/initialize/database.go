package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/xinliangnote/go-gin-api/cmd/config"
	"github.com/xinliangnote/go-gin-api/cmd/global"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitErpDB(ctx context.Context) {
	global.Logger.Info(ctx, "InitErpDB...")
	mysqlInfo := global.ServerConfig.ErpMysqlInfo
	var err error
	global.DB, err = newDBEngine(&mysqlInfo)
	if err != nil {
		global.Logger.Errorf(ctx, "InitErpDB 错误：%s", err.Error())
		panic(fmt.Errorf("InitErpDB 错误：%s", err.Error()))
	}
}

func newDBEngine(dbConfig *config.DatabaseConfig) (*gorm.DB, error) {
	config := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix: "t_",   // table name prefix, table for `User` would be `t_users`
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			// NoLowerCase:   true, // skip the snake_casing of names
			// NameReplacer: strings.NewReplacer("CID", "Cid"), // use name replacer to change struct/field name before convert it to db name
		},
		Logger: logger.Default.LogMode(logger.Silent), // 设置为静默模式，不打印日志
	}

	if global.ServerConfig.AppInfo.RunMode == "debug" {
		config.Logger = logger.Default.LogMode(logger.Info)
	}

	s := "%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local"
	db, err := gorm.Open(mysql.Open(fmt.Sprintf(s,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Name,
		dbConfig.Charset,
		dbConfig.ParseTime,
	)), config,
	)
	if err != nil {
		return nil, err
	}

	// 获取通用数据库对象 sql.DB ，然后使用其提供的功能
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}
