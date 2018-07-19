package main

import (
	"log"
	"time"

	"github.com/anonutopia/gowaves"
)

type EthereumMonitor struct {
}

func (em *EthereumMonitor) start() {
	go func() {
		for {
			em.checkAddresses()
			time.Sleep(time.Second * 5)
		}
	}()
}

func (em *EthereumMonitor) checkAddresses() {
	var users []*User
	db.Where("ether_balance_new > 0").Find(&users)

	for _, u := range users {
		em.sendEther(u)
	}
}

func (em *EthereumMonitor) sendEther(u *User) {
	atr := &gowaves.AssetsTransferRequest{
		Amount:    u.EtherBalanceNew,
		AssetID:   "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53",
		Fee:       100000,
		Recipient: u.Address,
		Sender:    conf.NodeAddress,
	}

	_, err := wnc.AssetsTransfer(atr)

	if err != nil {
		log.Printf("[EthereumMonitor.sendEther] error assets transfer: %s", err)
	} else {
		log.Printf("Sent ETH: %s => %d", u.Address, u.EtherBalanceNew)
		u.EtherBalanceProcessed += u.EtherBalanceNew
		u.EtherBalanceNew = 0
		db.Save(u)
	}
}

func initEthMonitor() *EthereumMonitor {
	em := &EthereumMonitor{}
	em.start()
	return em
}
