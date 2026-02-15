package service

import "errors"

var ErrWalletNotFound = errors.New("wallet not found")
var ErrNotEnoughMoney = errors.New("not enough money")
