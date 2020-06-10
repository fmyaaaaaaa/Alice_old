package main

import (
	"flag"
	"fmt"
	"github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/config"
	database2 "github.com/fmyaaaaaaa/Alice/alice-trading/infrastructure/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database"
	"github.com/fmyaaaaaaa/Alice/alice-trading/interfaces/database/instruments"
	"github.com/fmyaaaaaaa/Alice/alice-trading/usecase"
)

func main() {
	flag.Parse()
	config.InitInstance("./infrastructure/config/env/", flag.Arg(0))
	fmt.Println("application started :", flag.Arg(0))

	// TODO:DBアクセスの暫定実装のため、不要になり次第修正する
	i := usecase.InstrumentsInteractor{
		DB:          &database.DBRepository{DB: database2.NewDB()},
		Instruments: &instruments.Repository{},
	}

	res, _ := i.Get(1)
	fmt.Println(res.Name)
}
