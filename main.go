package main

import (
	"github.com/anonutopia/gowaves"
	"github.com/jinzhu/gorm"
	"gopkg.in/telegram-bot-api.v4"
)

var conf *Config

var wnc *gowaves.WavesNodeClient

var db *gorm.DB

var pc *PriceClient

var bm *BitcoinMonitor

var em *EthereumMonitor

var bot *tgbotapi.BotAPI

func main() {
	conf = initConfig()

	wnc = initWaves()

	db = initDb()

	pc = initPriceClient()

	bm = initBtcMonitor()

	em = initEthMonitor()

	bot = initBot()

	initMonitor()
}
