package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
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

	toAccDef := accounts.Account{
		Address: common.HexToAddress(to),
	}

	signAcc, err := eg.keystore.Find(fromAccDef)
	if err != nil {
		log.Printf("account keystore find error: %s", err)
		return err
	}

	errUnlock := eg.keystore.Unlock(signAcc, conf.BtcMasterKey)
	if errUnlock != nil {
		fmt.Printf("account unlock error: %s", err)
		return errUnlock
	}

	// Construct the transaction
	tx := types.NewTransaction(
		0x0,
		toAccDef.Address,
		big.NewInt(int64(amount)*(10^10)),
		uint64(21000),
		big.NewInt(50),
		[]byte("forward"))

	signedTx, errSign := eg.keystore.SignTx(signAcc, tx, big.NewInt(1))
	if errSign != nil {
		fmt.Printf("tx sign error: %s", err)
		return errSign
	}

	client, errDial := ethclient.Dial("http://localhost:8545")
	if errDial != nil {
		fmt.Printf("Dial error: %s", err)
		return errDial
	}

	txErr := client.SendTransaction(context.Background(), signedTx)
	if txErr != nil {
		fmt.Printf("send tx error: %s", txErr)
		return txErr
	}

	return nil
}

func initEthGen() *EthereumGenerator {
	eg := &EthereumGenerator{}
	eg.keystore = keystore.NewKeyStore("/home/kriptokuna/wallet/wallets", keystore.StandardScryptN, keystore.StandardScryptP)
	return eg
}
