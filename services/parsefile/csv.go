package parsefile

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

const headerPrefix = "header"

type csvDetails struct {
	data    []byte
	headers string
}

func NewCSV(headers string, data []byte) *csvDetails {
	return &csvDetails{data: data, headers: headers}
}

func (c *csvDetails) Read() ([]map[string]string, error) {

	csvReader := csv.NewReader(bytes.NewReader(c.data))

	records, err := csvReader.ReadAll()
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("error reading CSV")
		return nil, err
	}

	headerList := c.getHeadersList(len(records))
	recordMapsList := make([]map[string]string, 0)
	for _, record := range records {
		recordMap := make(map[string]string, 0)

		for i, recordEntry := range record {
			recordMap[headerList[i]] = recordEntry
		}
		recordMapsList = append(recordMapsList, recordMap)

	}
	return recordMapsList, nil
}

func (c *csvDetails) getHeadersList(numOfRecords int) []string {

	headerList := []string{}
	if c.headers == "" {
		for i := 0; i < numOfRecords; i++ {
			headerList = append(headerList, headerPrefix+strconv.Itoa(i))
		}
		return headerList
	}

	return strings.Split(c.headers, ",")
}
