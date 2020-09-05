package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type s3FileStoreSuite struct {
	suite.Suite
	store store
}

func (suite *s3FileStoreSuite) SetupSuite() {
	store := NewFileS3Store("https://focusmobile-interview-materials.s3.eu-west-3.amazonaws.com/Cheap.Stocks.Internationalization.Currencies.csv", "./modified.lock", "currency.cache")

	suite.store = store
}

func (suite *s3FileStoreSuite) TestGetCurrencies() {
	currencies, err := suite.store.getCurrencies()
	fmt.Print(currencies)

	suite.NoError(err)
	suite.NotNil(currencies)
	suite.NotZero(len(currencies))
}

func TestS3FileStoreSuite(t *testing.T) {
	suite.Run(t, new(s3FileStoreSuite))
}
