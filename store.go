package main

type store interface {
	getCurrencies() ([]currency, error)
}

type currency struct {
	Country     string
	Description string
	Code        string
}
