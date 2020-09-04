package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const lastModifiedFormat = "Mon, 02 Jan 2006 15:04:05 MST" //http times are RFC2616

var versionFile string = ""
var currencyFile string = ""

func readHTTPS3File(url string, ctx context.Context) ([]byte, error) {

	//url := currencyFile // a la global. this function is used in one place. ie behind an interface. it is mockable
	req, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return nil, errors.Wrap(err, "Error creating request")
	}

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		if err == context.Canceled {
			return nil, errors.Wrap(err, "Request timeout")
		}
		return nil, errors.Wrap(err, "Error with http request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(err, "Got invalid response from the http request")
	}

	var data *bytes.Buffer //zero value of this type is a valid Buffer. awesome!
	_, err = io.Copy(data, resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Error reading data from http response")
	}

	return data.Bytes(), nil
}

func preFlightRequest(url string, ctx context.Context) (http.Header, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)

	if err != nil {
		return nil, errors.Wrap(err, "Error creating preflight request")
	}

	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, errors.Wrap(err, "Error with getting preflight response")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(err, "Error with preflight response status. Not ok")
	}

	return resp.Header, nil

}

//This function uses a combination of Last-Modified  and Content-Length header
//I did not use the Etag header as aws is not consistent with what the hash represents
//sometimes it is a hash of the file contents or some metadata we don't care about.

//This function is stateful, it reads a modified.lock file to detect file changes
func s3FileNotModified(versionFilepath string, headers http.Header) bool {

	var ret bool = false
	modTime := headers.Get("Last-Modified")

	remoteTime, _ := time.Parse(lastModifiedFormat, modTime) // I don't know what to do with errs here. but shd panic

	data, _ := ioutil.ReadFile(versionFilepath)

	sdata := strings.Split(string(data), "|")

	cachedTime, _ := time.Parse(lastModifiedFormat, sdata[0])

	remoteLen := headers.Get("Content-Length")
	cachedLen := sdata[1]

	if strings.TrimSpace(remoteLen) == strings.TrimSpace(cachedLen) && remoteTime.Equal(cachedTime) {
		ret = true
	}

	return ret
}

func WriteCSV(currencies []currency) error {

	records := make([][]string, 1)

	for _, value := range currencies {

		record := structToSlice(value)
		records = append(records, record)
	}

	fobj, err := os.OpenFile(currencyFile, os.O_RDWR|os.O_TRUNC, 0644)

	if err != nil {
		return errors.Wrap(err, "Error opening currency file")
	}
	csvEncoder := csv.NewWriter(fobj)

	err = csvEncoder.WriteAll(records)

	if err != nil {
		return errors.Wrap(err, "Error encoding the data to csv")
	}

	csvEncoder.Flush() // Writes are buffered by this encoder. you need to flush to the writer

	return nil
}

func readCSV(r io.Reader) ([]currency, error) {
	currencies := [][]string{}
	csvReader := csv.NewReader(r)

	currencies, err := csvReader.ReadAll()

	if err != nil {
		return nil, errors.Wrap(err, "Error reading and decoding csv file")
	}

	return sliceToStruct(currencies), nil
}
