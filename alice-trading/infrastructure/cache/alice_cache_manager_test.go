package cache

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/domain/enum"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type TestModel struct {
	Name string
}

func TestGetCacheManager(t *testing.T) {
	// Config-test
	dummyConf := []string{"", "./../config/env/"}
	config.InitInstance("test", dummyConf)

	// FirstTime
	cacheManagerFirst := GetCacheManager()
	assert.Equal(t, reflect.TypeOf(AliceCacheManager{}).String(), reflect.TypeOf(cacheManagerFirst).String())

	// SecondTime
	cacheManagerSecond := GetCacheManager()
	assert.Equal(t, cacheManagerFirst, cacheManagerSecond)
}

func TestAliceCacheManager_Set(t *testing.T) {
	model1 := TestModel{Name: "Sample1"}
	cacheManager := GetCacheManager()
	cacheManager.Set("test", model1, enum.DEFAULT)
	data := cacheManager.Get("test")
	assert.Equal(t, model1, data)

	model2 := TestModel{Name: "Sample2"}
	cacheManager.Set("test", model2, enum.DEFAULT)
	data = cacheManager.Get("test")
	assert.NotEqual(t, model1, data)
	assert.Equal(t, model2, data)

	model3 := []TestModel{{Name: "Sample3"}, {Name: "Sample4"}}
	cacheManager.Set("array", model3, enum.DEFAULT)
	data = cacheManager.Get("array")
	targetArray := data.([]TestModel)
	assert.Equal(t, 2, len(targetArray))
	assert.Equal(t, "Sample4", targetArray[1].Name)
}
