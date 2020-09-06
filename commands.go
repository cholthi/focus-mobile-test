package main

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

const INVALIDINPUT = 400
const SYSTEMERROR = 500

func CurrencyExists(store store) func(*cli.Context) error {
	return func(c *cli.Context) error {
		currency := c.Args().Get(0)
		if ok := isISO4217(currency); !ok {
			return cli.Exit("The argument passed is not valid ISO 4217 currency code", INVALIDINPUT)
		}
		currencies, err := store.getCurrencies()

		if err != nil {
			return cli.Exit("Error: "+err.Error(), SYSTEMERROR)
		}
		var supported bool = false
		for i := range currencies {
			if currencies[i].Code == strings.ToUpper(currency) {
				supported = true
			}
		}

		if supported {
			fmt.Printf("The currency %s is supported \n", currency)
			return nil
		}

		fmt.Printf("The currency %s is not supported\n", currency)
		return nil

	}
}
