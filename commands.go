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
		cur := c.Args().Get(0)
		if ok := isISO4217(cur); !ok {
			return cli.Exit("The argument passed is not valid ISO 4217 currency", INVALIDINPUT)
		}
		currencies, err := store.getCurrencies()

		if err != nil {
			return cli.Exit("Error: "+err.Error(), SYSTEMERROR)
		}
		var supported bool = false
		for i := range currencies {
			if currencies[i].Code == strings.ToUpper(cur) {
				supported = true
			}
		}

		if supported {
			fmt.Printf("The currency %s is supported \n", cur)
			return nil
		}

		fmt.Printf("The currency %s is not supported\n", cur)
		return nil

	}
}
