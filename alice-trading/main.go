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
	if config.InitInstance(flag.Arg(0), flag.Args()) {
		fmt.Println("application started :", flag.Arg(0))
		fmt.Println("param :", flag.Args())
		// TODO:DBアクセスの暫定実装のため、不要になり次第修正する
		i := usecase.InstrumentsInteractor{
			DB:          &database.DBRepository{DB: database2.NewDB()},
			Instruments: &instruments.Repository{},
		}
		fmt.Println(i.Get(1))
	} else {
		panic("failed application initialize")
	}
}
