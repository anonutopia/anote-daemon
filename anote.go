package main

import "log"

const satInBtc = uint64(100000000)

const priceFactorLimit = uint64(0.0001 * float64(satInBtc))

type Anote struct {
	Price           uint64
	PriceFactor     uint64
	TierPrice       uint64
	TierPriceFactor uint64
}

func (a *Anote) issueAmount(investment int, assetID string) int {
	p, err := pc.DoRequest()
	amount := int(0)
	if err == nil {
		var cryptoPrice float64

		if len(assetID) == 0 {
			cryptoPrice = p.WAVES
		} else if assetID == "7xHHNP8h6FrbP5jYZunYWgGn2KFSBiWcVaZWe644crjs" {
			cryptoPrice = p.BTC
		} else if assetID == "4fJ42MSLPXk9zwjfCdzXdUDAH8zQFCBdBz4sFSWZZY53" {
			cryptoPrice = p.ETH
		} else {
			return amount
		}

		// amount := int(float64(t.Amount) / cryptoPrice / anote.Price)
		for investment > 0 {
			tierAmount := uint64(float64(investment) / cryptoPrice / float64(a.Price) * float64(satInBtc))

			log.Printf("tierAmount: %s", tierAmount)

			if tierAmount > a.TierPrice {
				tierAmount = a.TierPrice
			}

			amount += int(tierAmount)
			investment -= int(tierAmount * a.Price / satInBtc)

			a.TierPrice -= tierAmount
			a.TierPriceFactor -= tierAmount

			if a.TierPrice == 0 {
				a.TierPrice = 1000 * satInBtc
				a.Price = a.Price + a.PriceFactor
				a.saveState()
			}

			if a.TierPriceFactor == 0 {
				a.TierPriceFactor = 1000000 * satInBtc
				if a.PriceFactor > priceFactorLimit {
					a.PriceFactor = a.PriceFactor - priceFactorLimit
					a.saveState()
				}
			}
		}
	} else {
		log.Printf("[Anote.issueAmount] error pc.DoRequest: %s", err)
	}

	return amount
}

func (a *Anote) saveState() {
	ksip := &KeyValue{Key: "anotePrice"}
	db.FirstOrCreate(ksip)
	ksip.Value = a.Price
	db.Save(ksip)

	ksipf := &KeyValue{Key: "anotePriceFactor"}
	db.FirstOrCreate(ksipf)
	ksipf.Value = a.PriceFactor
	db.Save(ksipf)

	ksitp := &KeyValue{Key: "anoteTierPrice"}
	db.FirstOrCreate(ksitp)
	ksitp.Value = a.TierPrice
	db.Save(ksitp)

	ksitpf := &KeyValue{Key: "anoteTierPriceFactor"}
	db.FirstOrCreate(ksitpf)
	ksitpf.Value = a.TierPriceFactor
	db.Save(ksitpf)
}

func (a *Anote) loadState() {
	ksip := &KeyValue{Key: "anotePrice"}
	db.FirstOrCreate(ksip, ksip)

	if ksip.Value > 0 {
		a.Price = ksip.Value
	} else {
		ksip.Value = a.Price
		db.Save(ksip)
	}

	ksipf := &KeyValue{Key: "anotePriceFactor"}
	db.FirstOrCreate(ksipf, ksipf)

	if ksipf.Value > 0 {
		a.PriceFactor = ksipf.Value
	} else {
		ksipf.Value = a.PriceFactor
		db.Save(ksipf)
	}

	ksitp := &KeyValue{Key: "anoteTierPrice"}
	db.FirstOrCreate(ksitp, ksitp)

	if ksitp.Value > 0 {
		a.TierPrice = ksitp.Value
	} else {
		ksitp.Value = a.TierPrice
		db.Save(ksitp)
	}

	ksitpf := &KeyValue{Key: "anoteTierPriceFactor"}
	db.FirstOrCreate(ksitpf, ksitpf)

	if ksitpf.Value > 0 {
		a.TierPriceFactor = ksitpf.Value
	} else {
		ksitpf.Value = a.TierPriceFactor
		db.Save(ksitpf)
	}
}

func initAnote() *Anote {
	anote := &Anote{Price: uint64(0.01 * float64(satInBtc)), PriceFactor: uint64(0.0021 * float64(satInBtc)), TierPrice: 1000 * satInBtc, TierPriceFactor: 1000000 * satInBtc}

	anote.loadState()

	return anote
}
