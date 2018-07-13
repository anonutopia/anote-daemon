package main

import (
	"log"
	"time"

	"github.com/anonutopia/gowaves"
)

type WavesMonitor struct {
	StartedTime int64
}

func (wm *WavesMonitor) start() {
	wm.StartedTime = time.Now().Unix() * 1000
	for {
		// todo - make sure that everything is ok with 10 here
		pages, err := wnc.TransactionsAddressLimit("3PDb1ULFjazuzPeWkF2vqd1nomKh4ctq9y2", 100)
		if err != nil {
			log.Println(err)
		}
		if len(pages) > 0 {
			for _, t := range pages[0] {
				wm.checkTransaction(&t)
			}
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
	if t.Type == 4 && t.Timestamp >= wm.StartedTime {
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
					log.Printf("Sent ANO: %s => %d", t.Sender, amount)
				}

				user := &User{Address: t.Sender}
				db.First(user, user)
				if len(user.Referral) > 0 {
					referral := &User{Address: user.Referral}
					db.First(referral, referral)
					splitToHolders := t.Amount / 2
					splitToHolders -= (t.Amount / 5)
					if referral.ID != 0 {
						newProfit := uint64(t.Amount / 5)
						referral.ReferralProfitWav += newProfit
						referral.ReferralProfitWavTotal += newProfit
						db.Save(referral)
						splitToHolders -= (t.Amount / 5)
					}
					wm.splitToHolders(splitToHolders)
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

func (wm *WavesMonitor) splitToHolders(amount int) {
	// log.Println(amount)
	// ad, err := wnc.AssetsDistribution("4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	// log.Println(err)
	// log.Println(ad)
	// ts, err := wm.totalSupply()
	// if err == nil {
	// 	log.Println(ts)
	// }
	ad, err := wnc.AssetsDistribution("4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		log.Println(err)
	} else {
		itemsMap := ad.(map[string]interface{})
		for k, _ := range itemsMap {
			stake, err := wm.calculateStake(k)
			if err == nil {
				user := &User{Address: k}
				db.First(user, user)
				if user.ID != 0 {
					amountUser := uint64(float64(amount) * stake)
					user.ProfitWav += amountUser
					db.Save(user)
				}
			}
		}
	}
}

func (wm *WavesMonitor) totalSupply() (uint64, error) {
	supply := uint64(0)
	ad, err := wnc.AssetsDistribution("4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		return 0, err
	}
	itemsMap := ad.(map[string]interface{})
	for k, a := range itemsMap {
		if k != "3PDb1ULFjazuzPeWkF2vqd1nomKh4ctq9y2" {
			supply += uint64(a.(float64))
		}
	}
	return supply, err
}

func (wm *WavesMonitor) getBalance(address string) (uint64, error) {
	ad, err := wnc.AssetsDistribution("4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		return 0, err
	}
	itemsMap := ad.(map[string]interface{})
	for k, a := range itemsMap {
		if k != address {
			return uint64(a.(float64)), nil
		}
	}
	return uint64(0), err
}

func (wm *WavesMonitor) calculateStake(address string) (float64, error) {
	if address == "3PDb1ULFjazuzPeWkF2vqd1nomKh4ctq9y2" {
		return 0, nil
	}
	b, err := wm.getBalance(address)
	if err != nil {
		return 0, err
	}
	ts, err := wm.totalSupply()
	if err != nil {
		return 0, err
	}
	return float64(b) / float64(ts), nil
}

func initMonitor() {
	wm := &WavesMonitor{}
	wm.start()
}
