package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/anonutopia/gowaves"
)

type BitcoinMonitor struct {
}

func (bm *BitcoinMonitor) start() {
	go func() {
		for {
			bm.checkAddresses()
			time.Sleep(time.Second * 5)
		}
	}()
}

func (bm *BitcoinMonitor) checkAddresses() {
	var users []*User
	db.Where("bitcoin_balance_new > 0").Find(&users)

	for _, u := range users {
		bm.sendBitcoin(u)
	}
}

func (bm *BitcoinMonitor) sendBitcoin(u *User) {
	atr := &gowaves.AssetsTransferRequest{
		Amount:    u.BitcoinBalanceNew,
		AssetID:   "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs",
		Fee:       100000,
		Recipient: u.Address,
		Sender:    conf.NodeAddress,
	}

	_, err := wnc.AssetsTransfer(atr)

	if err != nil {
		log.Printf("[BitcoinMonitor.sendBitcoin] error assets transfer: %s", err)
	} else {
		log.Printf("Sent BTC: %s => %d", u.Address, u.BitcoinBalanceNew)
		u.BitcoinBalanceProcessed += u.BitcoinBalanceNew
		u.BitcoinBalanceNew = 0
		db.Save(u)
	}
}

func initBtcMonitor() *BitcoinMonitor {
	bm := &BitcoinMonitor{}
	bm.start()
	return bm
}

type BitcoinGenerator struct {
}

func (bg *BitcoinGenerator) sendBitcoin(address string, amount float64) error {
	log.Println(fmt.Sprintf("%.8f", amount))
	cmd := exec.Command("/usr/local/bin/electrum", "payto", address, fmt.Sprintf("%.8f", amount), "|", "/usr/local/bin/electrum", "broadcast -")
	cmd.Env = append(os.Environ(), "HOME=/home/kriptokuna")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println("Error in BitcoinGenerator.sendBitcoin: " + string(stderr.Bytes()))
		return err
	}
	return nil
}

func initBtcGen() *BitcoinGenerator {
	bg := &BitcoinGenerator{}
	return bg
}
