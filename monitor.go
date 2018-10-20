package main

import (
	"fmt"
	"log"
	"strings"
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
		pages, err := wnc.TransactionsAddressLimit(conf.NodeAddress, 100)
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
	if t.Type == 4 && t.Timestamp >= wm.StartedTime && len(t.Attachment) == 0 && t.Sender != conf.NodeAddress && t.Recipient == conf.NodeAddress {
		if len(t.AssetID) == 0 || t.AssetID == "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs" || t.AssetID == "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53" {

			amount := anote.issueAmount(t.Amount, t.AssetID)

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

			splitToFunders := t.Amount / 2
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
					splitToFunders -= (t.Amount / 10)
				}
			}

			wm.splitToFunders(splitToFunders, t.AssetID)

			wm.addToBudget(splitToFunders, t.AssetID)
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
						Amount:    int(btcAmount),
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
						Amount:    int(ethAmount),
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
		} else if strings.HasPrefix(string(dcd), "forwardbtc=") {
			user := &User{Address: t.Sender}
			db.First(user, user)
			if user.ID != 0 {
				if t.Amount > 50000 {
					amount := t.Amount - 50000
					err := bg.sendBitcoin(strings.Replace(string(dcd), "forwardbtc=", "", 1), float64(amount)/float64(satInBtc))
					if err != nil {
						log.Printf("Error in bg.sendBitcoin: %s", err)
					} else {
						user.BitcoinBalanceProcessed -= t.Amount
						db.Save(user)

						anote.GatewayProfitBtc += 50000
						anote.saveState()
					}
				}
			}
		} else if strings.HasPrefix(string(dcd), "forwardeth=") {
			if t.Amount > 100000 {
				user := &User{Address: t.Sender}
				db.First(user, user)
				if user.ID != 0 {
					ua := &UsedAddress{}
					db.Where("balance >= ?", t.Amount).First(ua)
					if ua.ID != 0 {
						amount := t.Amount - 100000
						err := eg.sendEther(ua.Address, strings.Replace(string(dcd), "forwardeth=", "", 1), float64(amount)/float64(satInBtc))
						if err != nil {
							log.Printf("Error in eg.sendEther: %s", err)
						} else {
							ua.Balance -= uint64(t.Amount)
							db.Save(ua)
							anote.GatewayProfitEth += 100000
							anote.saveState()
						}
					}
				}
			}
		} else {
			log.Printf("[WavesMonitor.processTransaction] %s %.8f", dcd, float64(t.Amount)/float64(satInBtc))
			msg := tgbotapi.NewMessage(-1001325718529, string(dcd))
			bot.Send(msg)
		}
	}

	tr.Processed = true
	db.Save(tr)
}

func (wm *WavesMonitor) splitToFunders(amount int, assetID string) {
	funder := &Badge{Name: "funder"}
	db.Preload("Users").First(funder)

	for _, user := range funder.Users {
		stake, err := wm.calculateStake(user.Address)
		if err == nil {
			log.Printf("Stake for address %s => %2f", user.Address, stake)

			amountUser := uint64(float64(amount) * stake)

			if len(assetID) == 0 {
				user.ProfitWav += amountUser
				user.ProfitWavTotal += amountUser
			} else if assetID == "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs" {
				user.ProfitBtc += amountUser
				user.ProfitBtcTotal += amountUser
			} else if assetID == "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53" {
				user.ProfitEth += amountUser
				user.ProfitEthTotal += amountUser
			}

			db.Save(user)
		} else {
			log.Printf("error in calculateStake: %s", err)
		}
	}
}

func (wm *WavesMonitor) addToBudget(amount int, assetID string) {
	if len(assetID) == 0 {
		anote.BudgetWav += uint64(amount)
	} else if assetID == "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs" {
		anote.BudgetBtc += uint64(amount)
	} else if assetID == "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53" {
		anote.BudgetEth += uint64(amount)
	}

	anote.saveState()
}

func (wm *WavesMonitor) totalSupply() (uint64, error) {
	supply := uint64(0)

	funder := &Badge{Name: "funder"}
	db.Preload("Users").First(funder)

	for _, user := range funder.Users {
		balance, _ := wm.getBalance(user.Address)
		supply += balance
	}

	return supply, nil
}

func (wm *WavesMonitor) getBalance(address string) (uint64, error) {
	abr, err := wnc.AssetsBalance(address, "4zbprK67hsa732oSGLB6HzE8Yfdj3BcTcehCeTA1G5Lf")
	if err != nil {
		return 0, err
	}
	return uint64(abr.Balance), err
}

func (wm *WavesMonitor) calculateStake(address string) (float64, error) {
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
