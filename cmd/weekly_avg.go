package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"time"

	"github.com/meghashyamc/calorie-insights/services/parsefile"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	headerPrefix         = "header"
	calorieCSVTimeLayout = "1/02/06 3:04 PM"
	reqDateLayout        = "2006-01-02"
)

var csvForWeeklyAvg *string

type weeklyAvgDataList struct {
	data []weeklyAvgData
}

type weeklyAvgData struct {
	weekStartDate string
	weekEndDate   string
	avgData       int
}

var weeklyAvgCmd = &cobra.Command{
	Use:   "weeklyavg",
	Short: "get weekly average calories from a calorie CSV",
	Long:  "get weekly average calories from a calorie CSV",
	Run:   getWeeklyAverage,
}

func getWeeklyAverage(cmd *cobra.Command, args []string) {

	csvToUse, err := getCSVForWeeklyAvg()
	if err != nil {
		return
	}
	fileBytes, err := ioutil.ReadFile(csvToUse)
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not read from passed filepath")
		return
	}

	newCSV := parsefile.NewCSV("", fileBytes)

	fileData, err := newCSV.Read()
	if err != nil {
		return
	}

	gottenWeeklyAvgData, err := getWeeklyAvgData(fileData)
	if err != nil {
		return
	}

	gottenWeeklyAvgData.print()

}

func getWeeklyAvgData(fileData []map[string]string) (*weeklyAvgDataList, error) {

	dateCaloriesMap, sortedDates, err := getDateCaloriesMap(fileData)
	if err != nil {
		return nil, err
	}

	result := &weeklyAvgDataList{}
	result.data, err = calculateWeeklyAvgData(sortedDates, dateCaloriesMap)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func calculateWeeklyAvgData(sortedDates []string, dateCaloriesMap map[string]int) ([]weeklyAvgData, error) {

	weeklyAvgDataList := []weeklyAvgData{}
	for i := 0; i < len(sortedDates); {

		weekStartDate := sortedDates[i]
		weekEndDate, err := getEndOfWeekGivenStart(sortedDates[i])
		if err != nil {
			return nil, err
		}
		thisWeeksAvgData := weeklyAvgData{weekStartDate: weekStartDate, weekEndDate: weekEndDate}
		var intervalSize int
		thisWeeksAvgData.avgData, intervalSize = getAvgDataForDates(weekEndDate, i, sortedDates, dateCaloriesMap)

		weeklyAvgDataList = append(weeklyAvgDataList, thisWeeksAvgData)
		i += intervalSize
	}

	return weeklyAvgDataList, nil

}

func getAvgDataForDates(endDate string, sortedDatesIndex int, sortedDates []string, dateCaloriesMap map[string]int) (int, int) {

	caloriesInInterval := 0
	intervalSize := 0

	for _, date := range sortedDates[sortedDatesIndex:] {

		if date > endDate {
			break
		}

		dateCalories, ok := dateCaloriesMap[date]
		if !ok {
			continue
		}
		caloriesInInterval += dateCalories
		intervalSize++
	}
	return caloriesInInterval / intervalSize, intervalSize

}

func getEndOfWeekGivenStart(startDate string) (string, error) {

	startDateAsTime, err := time.Parse(reqDateLayout, startDate)
	if err != nil {
		log.WithFields(log.Fields{"date_to_parse": startDate, "err": err.Error()}).Error("could not parse date string to time")
		return "", err
	}

	return startDateAsTime.Add(6 * 24 * time.Hour).Format(reqDateLayout), nil

}

func getDateCaloriesMap(fileData []map[string]string) (map[string]int, []string, error) {

	dateCaloriesMap := map[string]int{}
	dates := []string{}

	for _, entry := range fileData {
		date, calories, err := validateCaloriesEntry(entry)
		if err != nil {
			return nil, nil, err
		}
		existingCalories, ok := dateCaloriesMap[date]
		if !ok {
			dateCaloriesMap[date] = calories
			dates = append(dates, date)
			continue
		}
		dateCaloriesMap[date] = existingCalories + calories

	}

	sort.Slice(dates, func(i, j int) bool {
		return dates[i] < dates[j]
	})

	return dateCaloriesMap, dates, nil
}

func validateCaloriesEntry(entry map[string]string) (string, int, error) {
	dateHeader := headerPrefix + "0"
	calorieHeader := headerPrefix + "2"

	timeStr, ok := entry[dateHeader]
	if !ok {
		err := errors.New("could not find expected date value corresponding to header value after successfully parsing CSV")
		log.WithFields(log.Fields{"header_name": dateHeader, "entry": entry}).Error(err.Error())
		return "", 0, err
	}

	calorieStr, ok := entry[calorieHeader]
	if !ok {
		err := errors.New("could not find expected calorie value corresponding to header value after successfully parsing CSV")
		log.WithFields(log.Fields{"header_name": calorieHeader, "entry": entry}).Error(err.Error())
		return "", 0, err
	}

	parsedTime, err := time.Parse(calorieCSVTimeLayout, timeStr)

	if err != nil {
		log.WithFields(log.Fields{"err": err.Error(), "entry": entry, "expected_layout": calorieCSVTimeLayout}).Error("could not parse time in CSV")
		return "", 0, err
	}

	date := parsedTime.Format(reqDateLayout)

	calories, err := strconv.Atoi(calorieStr)
	if err != nil {
		log.WithFields(log.Fields{"calories_str": calorieStr}).Error("could not convert read calories from CSV to number")
		return "", 0, err
	}
	return date, calories, nil
}
func (w *weeklyAvgDataList) print() {
	for _, weekData := range w.data {
		fmt.Println(fmt.Sprintf("%+v", weekData))
	}

}

func getCSVForWeeklyAvg() (string, error) {

	if *csvForWeeklyAvg == "" {
		csvPath, err := getAlreadyAddedCSV()
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Error("getting weekly average requires a CSV file path argument or an already added CSV")
			return "", err
		}
		csvForWeeklyAvg = &csvPath
	}
	return *csvForWeeklyAvg, nil

}

func setupWeeklyAvgCmd() {
	csvForWeeklyAvg = weeklyAvgCmd.Flags().StringP("csv", "f", "", "--csv <CSV to use>")
}
