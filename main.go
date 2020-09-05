package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

const remotecurrencyFile string = "https://focusmobile-interview-materials.s3.eu-west-3.amazonaws.com/Cheap.Stocks.Internationalization.Currencies.csv"

func main() {
	var store store

	app := &cli.App{
		Name:    "Check company Supported Currency",
		Version: "v0.1",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Philip Chol Paul",
				Email: "cholthipaul@gmail.com",
			},
		},
		Copyright: "(c) 2020 Cheap Stocks Inc",
		Usage:     "AppName supported --currency `ISO 4217 Code`",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cacheFile, c",
				Value: "./currency.cache",
				Usage: "The `FILE` to cache contents of the remote currency url",
			},
			&cli.StringFlag{
				Name:  "versionFile, v",
				Value: "./modified.lock",
				Usage: "The `FILE` to keep track of remote file changes. It must be modified outside the app",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "supported",
				Aliases: []string{"s"},
				Usage:   "Determine if the ISO 4217 currency code is supported by the company",
				Action: func(c *cli.Context) error {
					store = NewFileS3Store(remotecurrencyFile, c.String("versionFile"), c.String("cacheFile"))
					action := CurrencyExists(store)

					return action(c)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
