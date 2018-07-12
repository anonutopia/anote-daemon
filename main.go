package main

import "github.com/anonutopia/gowaves"

var conf *Config

var wnc *gowaves.WavesNodeClient

func main() {
	conf = initConfig()

	wnc = initWaves()

	initMonitor()
}
