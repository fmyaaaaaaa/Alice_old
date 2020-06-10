package config

import (
	"github.com/BurntSushi/toml"
	"sync"
)

var instance *AliceConfig
var once sync.Once

type AliceConfig struct {
	Api Api
	DB  DB
}

type Api struct {
	Url         string `toml:"url"`
	AccountId   string `toml:"Account_id"`
	AccessToken string `toml:"Access_token"`
}

type DB struct {
	Host     string `toml:"host"`
	Port     string `toml:"port"`
	UserName string `toml:"User_name"`
	Password string `toml:"password"`
	DBName   string `toml:"Db_name"`
}

func NewConfig(path string, appMode string) (AliceConfig, error) {
	var conf AliceConfig
	confPath := path + appMode + ".toml"
	if _, err := toml.DecodeFile(confPath, &conf); err != nil {
		return conf, err
	}
	*instance = conf
	return *instance, nil
}

// AliceConfigの初期化を行います。
// 初期化処理が成功した場合にtrue、初期化処理の失敗または既に初期化済の場合にfalseを返却します。
func InitInstance(path string, appMode string) bool {
	res := false
	once.Do(func() {
		instance = &AliceConfig{}
		confPath := path + appMode + ".toml"
		if _, err := toml.DecodeFile(confPath, &instance); err == nil {
			res = true
		}
	})
	return res
}

func GetInstance() *AliceConfig {
	return instance
}
