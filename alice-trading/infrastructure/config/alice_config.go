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
	Rule     Rule
}

type (
	Api struct {
		Url         string `toml:"url"`
		AccountId   string `toml:"Account_id"`
		AccessToken string `toml:"Access_token"`
	}

	DB struct {
		Host     string `toml:"host"`
		Port     string `toml:"port"`
		UserName string `toml:"User_name"`
		Password string `toml:"password"`
		DBName   string `toml:"Db_name"`
	}

	Property struct {
		RiskTolerancePrice int     `toml:"Risk_tolerance_price"`
		OrderLot           int     `toml:"Order_lot"`
		ProfitGainPrice    float64 `toml:"Profit_gain_price"`
	}

	Rule struct {
		TradeRule string
		TradeType string
	}
)

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
		case "local", "test", "backTest":
			res = initLocalOrTest(appMode, instance, confParams)
			res = initLocalOrTest("common", instance, confParams)
		case "develop", "staging", "production":
			res = initSystemProperty(instance)
			res = initDevelopOrStagingOrProduction(instance, confParams)
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
	if len(confParams) != 11 {
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

	// Rule
	tradeRule := confParams[9]
	tradeType := confParams[10]

	conf := AliceConfig{Api: Api{Url: url, AccountId: accountID, AccessToken: accessToken},
		DB:       DB{Host: host, Port: port, UserName: userName, Password: password, DBName: dbName},
		Property: Property{RiskTolerancePrice: instance.Property.RiskTolerancePrice, OrderLot: instance.Property.OrderLot, ProfitGainPrice: instance.Property.ProfitGainPrice},
		Rule:     Rule{TradeRule: tradeRule, TradeType: tradeType}}
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
