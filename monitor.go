package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/mr-tron/base58/base58"
	"gopkg.in/telegram-bot-api.v4"
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
	if tr.Processed != true {
		wm.processTransaction(&tr, t)
	}
}

func (wm *WavesMonitor) processTransaction(tr *Transaction, t *gowaves.TransactionsAddressLimitResponse) {
	if t.Type == 4 && t.Timestamp >= wm.StartedTime && len(t.Attachment) == 0 && t.Sender != "3PDb1ULFjazuzPeWkF2vqd1nomKh4ctq9y2" && t.Recipient == "3PDb1ULFjazuzPeWkF2vqd1nomKh4ctq9y2" {
		if len(t.AssetID) == 0 || t.AssetID == "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs" || t.AssetID == "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53" {
			p, err := pc.DoRequest()
			if err == nil {
				var cryptoPrice float64
				var invType string
				if len(t.AssetID) == 0 {
					cryptoPrice = p.WAVES
					invType = "WAV"
				} else if t.AssetID == "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs" {
					cryptoPrice = p.BTC
					invType = "BTC"
				} else if t.AssetID == "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53" {
					cryptoPrice = p.ETH
					invType = "ETH"
				}

				amount := int(float64(t.Amount) / cryptoPrice / 0.01)

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

				splitToHolders := t.Amount / 2
				user := &User{Address: t.Sender}
				db.First(user, user)
				if len(user.Referral) > 0 {
					referral := &User{Address: user.Referral}
					db.First(referral, referral)
					if referral.ID != 0 {
						newProfit := uint64(t.Amount / 5)
						if len(t.AssetID) == 0 {
							referral.ReferralProfitWav += newProfit
							referral.ReferralProfitWavTotal += newProfit
						} else if t.AssetID == "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs" {
							referral.ReferralProfitBtc += newProfit
							referral.ReferralProfitBtcTotal += newProfit
						} else if t.AssetID == "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53" {
							referral.ReferralProfitEth += newProfit
							referral.ReferralProfitEthTotal += newProfit
						}
						db.Save(referral)
						splitToHolders -= (t.Amount / 5)
					}
				}
				wm.splitToHolders(splitToHolders, invType)
			} else {
				log.Printf("[WavesMonitor.processTransaction] error pc.DoRequest: %s", err)
			}
		}
	} else if len(t.Attachment) > 0 {
		dcd, err := base58.Decode(t.Attachment)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(dcd))

		if string(dcd) == "withdraw" {
			user := &User{Address: t.Sender}
			db.First(user, user)
			if user.ID != 0 {
				wavAmount := user.ProfitWav + user.ReferralProfitWav
				btcAmount := user.ProfitBtc + user.ReferralProfitBtc
				ethAmount := user.ProfitEth + user.ReferralProfitEth

				user.ProfitWav = 0
				user.ReferralProfitWav = 0
				user.ProfitBtc = 0
				user.ReferralProfitBtc = 0
				user.ProfitEth = 0
				user.ReferralProfitEth = 0

				db.Save(user)

				if wavAmount > 0 {
					atr := &gowaves.AssetsTransferRequest{
						Amount:    int(wavAmount),
						Fee:       100000,
						Recipient: t.Sender,
						Sender:    conf.NodeAddress,
					}

					_, err := wnc.AssetsTransfer(atr)

					if err != nil {
						log.Printf("[WavesMonitor.processTransaction] Error: %s", err)
						msg := tgbotapi.NewMessage(-1001325718529, fmt.Sprintf("%s", err))
						bot.Send(msg)
					}
				}

				if btcAmount > 0 {
					atr := &gowaves.AssetsTransferRequest{
						Amount:    int(wavAmount),
						AssetID:   "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs",
						Fee:       100000,
						Recipient: t.Sender,
						Sender:    conf.NodeAddress,
					}

					_, err := wnc.AssetsTransfer(atr)

					if err != nil {
						log.Printf("[WavesMonitor.processTransaction] Error: %s", err)
						msg := tgbotapi.NewMessage(-1001325718529, fmt.Sprintf("%s", err))
						bot.Send(msg)
					}
				}

				if ethAmount > 0 {
					atr := &gowaves.AssetsTransferRequest{
						Amount:    int(wavAmount),
						AssetID:   "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53",
						Fee:       100000,
						Recipient: t.Sender,
						Sender:    conf.NodeAddress,
					}

					_, err := wnc.AssetsTransfer(atr)

					if err != nil {
						log.Printf("[WavesMonitor.processTransaction] Error: %s", err)
						msg := tgbotapi.NewMessage(-1001325718529, fmt.Sprintf("%s", err))
						bot.Send(msg)
					}
				}
			}
		} else {
			log.Printf("[WavesMonitor.processTransaction] %s %.8f", dcd, float64(t.Amount)/100000000)
			msg := tgbotapi.NewMessage(-1001325718529, string(dcd))
			bot.Send(msg)
		}
	}

	tr.Processed = true
	db.Save(tr)
}

func (wm *WavesMonitor) calculateAmount(trType int64, amount int64) int64 {
	var amountSending int64
	amountSending = 0
	return amountSending
}

func (wm *WavesMonitor) splitToHolders(amount int, invType string) {
	ad, err := wnc.AssetsDistribution("4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		log.Println(err)
	} else {
		itemsMap := ad.(map[string]interface{})
		for k, _ := range itemsMap {
			stake, err := wm.calculateStake(k)
			log.Printf("Stake for address %s => %2f", k, stake)
			if err == nil {
				user := &User{Address: k}
				db.FirstOrCreate(user, user)
				if user.ID != 0 {
					amountUser := uint64(float64(amount) * stake)
					if invType == "WAV" {
						user.ProfitWav += amountUser
						user.ProfitWavTotal += amountUser
					} else if invType == "BTC" {
						user.ProfitBtc += amountUser
						user.ProfitBtcTotal += amountUser
					} else if invType == "ETH" {
						user.ProfitEth += amountUser
						user.ProfitEthTotal += amountUser
					}
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
	abr, err := wnc.AssetsBalance(address, "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		return 0, err
	}
	return uint64(abr.Balance), err
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
	// log.Printf("Total supply: %d", ts)
	return float64(b) / float64(ts), nil
}

func initMonitor() {
	wm := &WavesMonitor{}
	wm.start()
}
