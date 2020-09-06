package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

const timeoutSeconds = time.Second * 30

type s3FileStore struct {
	RemoteUrl   string
	LocalPath   string
	VersionFile string
}

func NewFileS3Store(currencyURL, versionFile, currencyCacheFile string) *s3FileStore {
	store := &s3FileStore{RemoteUrl: currencyURL, LocalPath: currencyCacheFile, VersionFile: versionFile}

	return store
}

func (s3 *s3FileStore) getCurrencies() ([]currency, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSeconds)
	defer cancel()

	headers, err := preFlightRequest(s3.RemoteUrl, ctx)

	if err != nil {
		return nil, errors.Unwrap(err)
	}

	NotModified, err := s3FileNotModified(s3.VersionFile, headers)
	if err != nil {
		return nil, err
	}

	if NotModified { //Fetch from the local cache. No network roundtrip
		fobj, err := os.OpenFile(s3.LocalPath, os.O_RDONLY, 0644)

		if err != nil {
			return nil, errors.Wrap(err, "Error opening csv cache file. Please supply the --CacheFile or -c option with a valid file path.")
		}
		return readCSV(fobj)
	}

	err = s3.UpdateVersionFile(headers) //Don't forget to update our own version control file ;-)
	if err != nil {
		return nil, err
	}

	data, err := readHTTPS3File(s3.RemoteUrl, ctx)

	if err != nil {
		return nil, err
	}

	currencies, err := unMarshalCSV(data)

	if err != nil {
		return nil, err
	}
	fobj, err := os.OpenFile(s3.LocalPath, os.O_RDWR|os.O_TRUNC, 0644)

	if err != nil {
		return nil, errors.Wrap(err, "Error opening currency file")
	}
	err = WriteCSV(fobj, currencies) //cache the file contents locally

	if err != nil {
		return nil, errors.Wrap(err, "Error writing through the slice of currency")
	}

	return currencies, nil

}

func (s3 *s3FileStore) UpdateVersionFile(headers http.Header) error {
	content := fmt.Sprintf("%s|%s", headers.Get("Last-Modified"), headers.Get("Content-Length"))

	err := ioutil.WriteFile(s3.VersionFile, []byte(content), 0644)

	if err != nil {
		return err
	}

	return nil

}
