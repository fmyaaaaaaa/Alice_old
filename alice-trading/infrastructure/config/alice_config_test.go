package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewConfig(t *testing.T) {
	var result = []struct {
		name     string
		appMode  string
		expected AliceConfig
	}{
		{
			name:    "test",
			appMode: "test",
			expected: AliceConfig{
				Api: Api{
					Url:         "https://api-fxtest.oanda.com",
					AccountId:   "111-222-33333333-444",
					AccessToken: "token-test",
				},
				DB: DB{
					Host:     "test-postgres",
					Port:     "5432",
					UserName: "alice",
					Password: "alice",
					DBName:   "test_db",
				},
			},
		},
	}

	for _, r := range result {
		t.Run(r.name, func(t *testing.T) {
			confDir := "src/alice/config/env/"
			res, err := NewConfig(confDir, r.appMode)
			if err == nil {
				// API
				assert.Equal(t, r.expected.Api.Url, res.Api.Url)
				assert.Equal(t, r.expected.Api.AccountId, res.Api.AccountId)
				assert.Equal(t, r.expected.Api.AccessToken, res.Api.AccessToken)
				// DB
				assert.Equal(t, r.expected.DB.Host, res.DB.Host)
				assert.Equal(t, r.expected.DB.Port, res.DB.Port)
				assert.Equal(t, r.expected.DB.UserName, res.DB.UserName)
				assert.Equal(t, r.expected.DB.Password, res.DB.Password)
				assert.Equal(t, r.expected.DB.DBName, res.DB.DBName)
			}
		})
	}
}

func TestInitInstance(t *testing.T) {
	dummyConf := []string{"", "./env/"}
	// First time
	res := InitInstance("test", dummyConf)
	assert.Equal(t, true, res)
	// Second time
	res = InitInstance("test", dummyConf)
	assert.Equal(t, false, res)
}

func TestGetInstance(t *testing.T) {
	dummyConf := []string{"", "./env/"}
	InitInstance("test", dummyConf)
	res := GetInstance()
	// API
	assert.Equal(t, "https://api-fxtest.oanda.com", res.Api.Url)
	assert.Equal(t, "111-222-33333333-444", res.Api.AccountId)
	assert.Equal(t, "token-test", res.Api.AccessToken)
	// DB
	assert.Equal(t, "test-postgres", res.DB.Host)
	assert.Equal(t, "5432", res.DB.Port)
	assert.Equal(t, "alice", res.DB.UserName)
	assert.Equal(t, "alice", res.DB.Password)
	assert.Equal(t, "test_db", res.DB.DBName)
}

//func TestInitStagingInstance(t *testing.T) {
//	ResetInstance()
//	dummyConf := []string{
//		"staging",
//		"https://api-fxtest.oanda.com",
//		"111-222-33333333-444",
//		"token-test",
//		"localhost",
//		"5432",
//		"test",
//		"test",
//		"test_db",
//	}
//	res := InitInstance("staging", dummyConf)
//	assert.Equal(t, true, res)
//
//	instance = GetInstance()
//	assert.Equal(t, "https://api-fxtest.oanda.com", instance.Api.Url)
//	assert.Equal(t, "test_db", instance.DB.DBName)
//}
