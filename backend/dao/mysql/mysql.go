package mysql

import (
	"fmt"
	"time"
	"web/models"
	"web/setting"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// 加载mysql
func Init(cfg *setting.MysqlConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err = gorm.Open(mysql.Open(dsn))
	if err != nil {
		zap.L().Error("connect DB failed, err: %v\n", zap.Error(err))
		return
	}

	// 自动迁移表结构（根据Model创建/更新表）
	db.AutoMigrate(&models.User{}, &models.Community{})

	// 获取底层sql.DB并设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("get DB failed, err: %v\n", zap.Error(err))
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	return
}

func Close() {
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("get DB failed, err: %v\n", zap.Error(err))
	}
	sqlDB.Close()
}
