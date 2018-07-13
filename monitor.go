package main

import (
	"log"
	"time"

	"github.com/anonutopia/gowaves"
)

type WavesMonitor struct {
}

func (wm *WavesMonitor) start() {
	for {
		pages, err := wnc.TransactionsAddressLimit("3PDb1ULFjazuzPeWkF2vqd1nomKh4ctq9y2", 10)
		if err != nil {
			log.Println(err)
		}
		for _, t := range pages[0] {
			wm.checkTransaction(&t)
		}
		time.Sleep(time.Second)
	}
}

func (wm *WavesMonitor) checkTransaction(t *gowaves.TransactionsAddressLimitResponse) {
	tr := Transaction{TxId: t.ID}
	db.FirstOrCreate(&tr, &tr)
	if tr.Processed != 1 {
		wm.processTransaction(&tr, t)
	}
}

func (wm *WavesMonitor) processTransaction(tr *Transaction, t *gowaves.TransactionsAddressLimitResponse) {
	if t.Type == 4 {
		if len(t.AssetID) == 0 {
			p, err := pc.DoRequest()
			if err == nil {
				amount := int(float64(t.Amount) * p.WAVES / 0.01)
				atr := &gowaves.AssetsTransferRequest{
					Amount:    amount,
					AssetID:   "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf",
					Fee:       100000,
					Recipient: t.Sender,
					Sender:    conf.NodeAddress,
				}

				_, err := wnc.AssetsTransfer(atr)
				if err != nil {
					log.Printf("[WavesMonitor.processTransation] error assets transfer: %s", err)
				} else {
					log.Printf("Sent ANO: %s:%d", t.Sender, amount)
				}
			}
		}
	}

	tr.Processed = 1
	db.Save(tr)
}

func (wm *WavesMonitor) calculateAmount(trType int64, amount int64) int64 {
	var amountSending int64
	amountSending = 0
	return amountSending
}

func initMonitor() {
	wm := &WavesMonitor{}
	wm.start()
}
