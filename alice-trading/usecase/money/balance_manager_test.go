package money

import (
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/api/oanda"
	database2 "github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var balanceManager BalanceManager

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	dummyConf := []string{"", "./../../infrastructure/config/env/"}
	config.InitInstance("test", dummyConf)
	DB := database.NewDB()
	balanceManager = BalanceManager{
		DB:                 &database2.DBRepository{DB: DB},
		PricesApi:          oanda.NewPricesApi(),
		BackTestResult:     &database2.BackTestResultRepository{},
		BalanceManagements: &database2.BalanceManagementsRepository{},
	}
}

func TestBalanceManager_CalculationPips(t *testing.T) {
	// リスク許容額
	riskTolerancePrice := -3000.0
	basePrice := 107.45
	result := balanceManager.CalculationPips(riskTolerancePrice, basePrice)
	assert.Equal(t, float64(27), result)
}
