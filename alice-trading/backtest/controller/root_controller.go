package controller

import (
	"log"
	"net/http"
)

func StartServer() {
	// 売買ルールのコントローラーを初期化する。
	captainAmerica := CaptainAmericaController{}
	captainAmerica.Initialize()

	http.HandleFunc("/captain.america/setup", captainAmerica.handleSetup)
	http.HandleFunc("/captain.america/trade.plan", captainAmerica.handleTradePlan)

	log.Println("BackTest Server Starting at 7070")
	err := http.ListenAndServe(":7070", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
