package main

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
)

const timeoutSeconds = 10

type s3FileStore struct {
	RemoteUrl   string
	LocalPath   string
	VersionFile string
}

func (s3 *s3FileStore) getCurrencies() ([]currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds))
	defer cancel()

	headers, err := preFlightRequest(s3.RemoteUrl, ctx)

	if err != nil {
		return nil, err
	}

	if NotModified := s3FileNotModified(s3.VersionFile, headers); NotModified {
		fobj, err := os.OpenFile(s3.LocalPath, os.O_RDONLY, 0644)

		if err != nil {
			return nil, errors.Wrap(err, "Error opening csv file")
		}
		return readCSV(fobj)
	}

	data, err := readHTTPS3File(s3.RemoteUrl, ctx)

	if err != nil {
		return nil, err
	}

	currencies, err := unMarshalCSV(data)

	if err != nil {
		return nil, err
	}

	err = WriteCSV(currencies) //cache the file contents locally

	if err != nil {
		return nil, errors.Wrap(err, "Error writing through the slice of currency")
	}

	return currencies, nil

}
