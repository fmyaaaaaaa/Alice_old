package oanda

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewInstrumentsApi(t *testing.T) {
	_, err := NewInstrumentsApi()
	assert.NoError(t, err)
}

// TODO:APIのユニットテストの実装方針に従い修正する
//func TestInstrumentsApi_getInstruments(t *testing.T) {
//	api, err := NewInstrumentsApi()
//	assert.NoError(t, err)
//	res, err  := api.GetInstruments(context.Background())
//	assert.NoError(t, err)
//	assert.Equal(t, 71, len(res.Instruments))
//}
//
//func TestInstrumentsApi_getInstrument(t *testing.T) {
//	target := "USD_JPY"
//	api, err := NewInstrumentsApi()
//	assert.NoError(t, err)
//	res, err := api.GetInstrument(context.Background(), target)
//	assert.NoError(t, err)
//	assert.Equal(t, target, res.Instruments[0].Name)
//}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

// TODO:テストはDummyサーバに接続するように修正する
func setup() {
	config.InitInstance("./../../../config/env/", "stg")
}
