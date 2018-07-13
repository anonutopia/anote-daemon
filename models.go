package main

import "github.com/jinzhu/gorm"

type Transaction struct {
	gorm.Model
	TxId      string `sql:"size:255"`
	Processed uint64
}
