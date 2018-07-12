package main

import (
	"log"
	"time"
)

type WavesMonitor struct {
}

func (wm *WavesMonitor) start() {
	for {
		pages, err := wnc.TransactionsAddressLimit("3PDb1ULFjazuzPeWkF2vqd1nomKh4ctq9y2", 10)
		if err != nil {
			log.Println(err)
		}
		log.Println(pages[0][1].ID)
		time.Sleep(time.Second)
	}
}

func initMonitor() {
	wm := &WavesMonitor{}
	wm.start()
}
