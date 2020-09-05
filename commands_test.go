package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

type fakeStore struct {
}

func (fs *fakeStore) getCurrencies() ([]currency, error) {
	return []currency{
		currency{"Algeria", "Algerian Dinar", "DZD"},
		currency{"Angola", "Angolan kwanza", "AOA"},
		currency{"Kenya", "Kenyan shilling", "KES"},
	}, nil
}

func run(t *testing.T, args []string) {
	fakestore := &fakeStore{}
	app := cli.NewApp()
	//app.Action = CurrencyExists(fakeStore)
	app.Commands = []*cli.Command{
		&cli.Command{
			Name:      "supported",
			Usage:     "supported `CUR`",
			UsageText: "dtermines if the currency is supported by system",
			Action:    CurrencyExists(fakestore),
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "cacheFile, c",
			Value: "./currency.cache",
		},
		&cli.StringFlag{
			Name:  "versionFile, v",
			Value: "./modified.lock",
		},
	}

	err := app.Run(args)
	require.NoError(t, err)
}

func TestCurrencyExist(t *testing.T) {
	args := os.Args[0:1] // Name of the program.
	args = append(args, "supported")
	args = append(args, "DZD")
	run(t, args)
}
