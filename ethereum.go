package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/anonutopia/gowaves"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
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
		logTelegram(fmt.Sprintf("[EthereumMonitor.sendEther] error assets transfer: %s", err))
	} else {
		log.Printf("Sent ETH: %s => %d", u.Address, u.EtherBalanceNew)
		logTelegram(fmt.Sprintf("Sent ETH: %s => %d", u.Address, u.EtherBalanceNew))
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
		logTelegram(fmt.Sprintf("account keystore find error: %s", err))
		return err
	}

	errUnlock := eg.keystore.Unlock(signAcc, conf.BtcMasterKey)
	if errUnlock != nil {
		fmt.Printf("account unlock error: %s", errUnlock)
		return errUnlock
	}

	keyJson, readErr := ioutil.ReadFile(strings.Replace(signAcc.URL.String(), "keystore://", "", 1))
	if readErr != nil {
		fmt.Printf("key json read error: %s", readErr)
		return readErr
	}

	keyWrapper, keyErr := keystore.DecryptKey(keyJson, conf.BtcMasterKey)
	if keyErr != nil {
		fmt.Printf("key decrypt error: %s", keyErr)
		return keyErr
	}

	signer := types.HomesteadSigner{}

	client, errDial := ethclient.Dial("http://localhost:8545")
	if errDial != nil {
		fmt.Printf("Dial error: %s", errDial)
		return errDial
	}

	nonce, _ := client.NonceAt(context.Background(), signAcc.Address, nil)
	amountInt := amount * math.Pow(10, 18)
	// Construct the transaction
	tx := types.NewTransaction(
		nonce,
		toAccDef.Address,
		big.NewInt(int64(amountInt)),
		uint64(50000),
		big.NewInt(5000000000),
		[]byte("forward"))

	// signedTx, errSign := eg.keystore.SignTx(signAcc, tx, big.NewInt(1))
	// if errSign != nil {
	// 	fmt.Printf("tx sign error: %s", errSign)
	// 	return errSign
	// }

	signature, signatureErr := crypto.Sign(signer.Hash(tx).Bytes(), keyWrapper.PrivateKey)
	if signatureErr != nil {
		fmt.Printf("signature create error: %s", signatureErr)
	}

	signedTx, signErr := tx.WithSignature(signer, signature)
	if signErr != nil {
		fmt.Printf("signer with signature error:", signErr)
		return signErr
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
