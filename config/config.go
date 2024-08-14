package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	// 服务器设置
	Server struct {
		// 服务器监听端口
		Port int `toml:"port"`
	} `toml:"server"`

	// 数据库设置
	Database struct {
		// 数据库主机地址
		Host string `toml:"host"`
		// 数据库端口
		Port int `toml:"port"`
		// 数据库用户名
		User string `toml:"user"`
		// 数据库密码
		Password string `toml:"password"`
		// 数据库名称
		DBName string `toml:"db_name"`
	} `toml:"database"`

	// 环境设置
	Env struct {
		// 环境类型 development, production
		Type string `toml:"type"`
	} `toml:"env"`
}

// 配置文件对象工厂函数
func NewConfig() (*Config, error) {
	// 读取配置文件
	file, err := os.ReadFile("./configuration.toml")
	if err != nil {
		return nil, err
	}

	// 解析配置文件
	config := new(Config)
	err = toml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
