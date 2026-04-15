package controller

import (
	"fmt"
	"web/dao/mysql"
	"web/dao/redis"
	"web/logger"
	"web/pkg/snowflake"
	"web/setting"

	"go.uber.org/zap"
)

func Init() {
	// 加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("init setting failed, err: %v\n", err)
		return
	}

	// 初始化日志
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	zap.L().Info("===============================================")
	zap.L().Info("logger init success...")

	// 初始化mysql
	err := mysql.Init(setting.Conf.MysqlConfig)

	if err != nil {
		zap.L().Error("init mysql failed, err: ", zap.Error(err))
		return
	}
	zap.L().Info("mysql init success...")

	// 初始化redis
	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		zap.L().Error("init redis failed, err: ", zap.Error(err))
		return
	}

	zap.L().Info("redis init success...")

	//初始化snowflake
	if err := snowflake.Init(setting.Conf.AppConfig); err != nil {
		zap.L().Error("init snowflake failed, err: ", zap.Error(err))
		return
	}
	zap.L().Info("snowflake init success...")

	//初始化gin框架内置的校验器使用的翻译器
	if err := InitTrans("zh"); err != nil {
		zap.L().Error("init trans failed, err: ", zap.Error(err))
		return
	}
	zap.L().Info("trans init success...")
}

func Close() {
	defer zap.L().Sync()
	defer mysql.Close()
}
