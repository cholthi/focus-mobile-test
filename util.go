package main

import "bytes"

func structToSlice(currency currency) []string {
	ret := []string{currency.Country, currency.Description, currency.Code}
	return ret
}

func sliceToStruct(currencies [][]string) []currency {
	ret := make([]currency, 1)
	for _, record := range currencies {
		cur := currency{Country: record[0], Description: record[1], Code: record[2]}
		ret = append(ret, cur)
	}

	return ret
}

func unMarshalCSV(data []byte) ([]currency, error) {

	var buf *bytes.Buffer = bytes.NewBuffer(data)

	return readCSV(buf)
}
