package main

import (
	"github.com/anonutopia/gowaves"
	"github.com/jinzhu/gorm"
)

var conf *Config

var wnc *gowaves.WavesNodeClient

var db *gorm.DB

var pc *PriceClient

var bm *BitcoinMonitor

func main() {
	conf = initConfig()

	wnc = initWaves()

	db = initDb()

	pc = initPriceClient()

	bm = initBtcMonitor()

	initMonitor()
}
