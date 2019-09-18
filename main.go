package main

import (
	"github.com/anonutopia/gowaves"
)

var conf *Config

var wnc *gowaves.WavesNodeClient

var pc *PriceClient

// var bot *tgbotapi.BotAPI

// var anote *Anote

func main() {
	conf = initConfig()

	wnc = initWaves()

	pc = initPriceClient()

	// bot = initBot()

	// anote = initAnote()

	initMonitor()
}
