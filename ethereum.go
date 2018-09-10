package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
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

type EthereumGenerator struct {
	keystore *keystore.KeyStore
}

func (eg *EthereumGenerator) sendEther(from string, to string, amount float64) error {
	fromAccDef := accounts.Account{
		Address: common.HexToAddress(from),
	}

	signAcc, err := eg.keystore.Find(fromAccDef)
	if err != nil {
		log.Printf("account keystore find error: %s", err)
		return err
	}
	fmt.Printf("account found: signAcc.addr=%s; signAcc.url=%s\n", signAcc.Address.String(), signAcc.URL)
	fmt.Println()

	return nil
}

func initEthGen() *EthereumGenerator {
	eg := &EthereumGenerator{}
	eg.keystore = keystore.NewKeyStore("/home/kriptokuna/wallet/wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	return eg
}
