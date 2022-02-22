package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/meghashyamc/calorie-insights/services/parsefile"
	log "github.com/sirupsen/logrus"
)

const (
	headerPrefix         = "header"
	calorieCSVTimeLayout = "1/02/06 3:04 PM"
	reqDateLayout        = "2006-01-02"
)

func getAlreadyAddedCSV() (string, error) {
	return "calorie_counter.csv", nil
}

func getCSVFileData(headers, csvPath string) ([]map[string]string, error) {
	fileBytes, err := ioutil.ReadFile(csvPath)
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not read from passed filepath")
		return nil, err
	}

	newCSV := parsefile.NewCSV(headers, fileBytes)

	fileData, err := newCSV.Read()
	if err != nil {
		return nil, err
	}

	return fileData, nil
}

func getCSV(csvParam *string) (string, error) {

	if csvParam == nil {
		err := errors.New("CSV file path parameter was unexpectedly nil")
		log.WithFields(log.Fields{"err": err.Error()}).Error("CSV file path paramter should not be nil")
		return "", err
	}

	if *csvParam == "" {
		csvPath, err := getAlreadyAddedCSV()
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Error("getting weekly average requires a CSV file path argument or an already added CSV")
			return "", err
		}
		csvParam = &csvPath
	}
	return *csvParam, nil

}
