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

var bg *BitcoinGenerator

var em *EthereumMonitor

var eg *EthereumGenerator

var bot *tgbotapi.BotAPI

var anote *Anote

func main() {
	conf = initConfig()

	wnc = initWaves()

	db = initDb()

	pc = initPriceClient()

	bm = initBtcMonitor()

	bg = initBtcGen()

	em = initEthMonitor()

	eg = initEthGen()

	bot = initBot()

	anote = initAnote()

	initMonitor()
}
