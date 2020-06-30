package config

import (
	"github.com/BurntSushi/toml"
	"sync"
)

var instance *AliceConfig
var once sync.Once

type AliceConfig struct {
	Api      Api
	DB       DB
	Property Property
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

type Property struct {
	RiskTolerancePrice int `toml:"Risk_tolerance_price"`
	OrderLot           int `toml:"Order_lot"`
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
func InitInstance(appMode string, confParams []string) bool {
	res := false
	once.Do(func() {
		instance = &AliceConfig{}
		switch appMode {
		case "local", "test":
			res = initLocalOrTest(appMode, instance, confParams)
			res = initLocalOrTest("common", instance, confParams)
			break
		case "develop", "staging", "production":
			res = initSystemProperty(instance)
			res = initDevelopOrStagingOrProduction(instance, confParams)
			break
		default:
			panic("something wrong happened : appMode")
		}
	})
	return res
}

// Local/Testモードで起動した場合の初期化を行います。
func initLocalOrTest(confName string, instance *AliceConfig, confParams []string) bool {
	confPath := confParams[1] + confName + ".toml"
	if _, err := toml.DecodeFile(confPath, &instance); err == nil {
		return true
	}
	return false
}

// Develop/Staging/Productionモードで起動した場合の初期化を行います。
func initDevelopOrStagingOrProduction(instance *AliceConfig, confParams []string) bool {
	if len(confParams) != 9 {
		return false
	}
	// API
	url := confParams[1]
	accountID := confParams[2]
	accessToken := confParams[3]

	// DB
	host := confParams[4]
	port := confParams[5]
	userName := confParams[6]
	password := confParams[7]
	dbName := confParams[8]

	conf := AliceConfig{Api: Api{Url: url, AccountId: accountID, AccessToken: accessToken},
		DB:       DB{Host: host, Port: port, UserName: userName, Password: password, DBName: dbName},
		Property: Property{RiskTolerancePrice: instance.Property.RiskTolerancePrice, OrderLot: instance.Property.OrderLot}}
	*instance = conf
	return true
}

// システムプロパティを初期化します。
func initSystemProperty(instance *AliceConfig) bool {
	confPath := "./infrastructure/config/env/common.toml"
	if _, err := toml.DecodeFile(confPath, &instance); err == nil {
		return true
	}
	return false
}

func GetInstance() *AliceConfig {
	return instance
}
