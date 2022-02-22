package cmd

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type calorieSourceDataList struct {
	data []calorieSourceData
}

type calorieSourceData struct {
	source string
	data   int
}

var (
	csvForCalorieSources *string
	calorieSourceTags    *string
)

var calorieSourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "get percentages of calories from various sources",
	Long:  "get percentages of calories from various sources",
	Run:   getCalorieSourceDetails,
}

func getCalorieSourceDetails(cmd *cobra.Command, args []string) {

	calorieSources, valid := validateCalorieSourcesTags(*calorieSourceTags)
	if !valid {
		return
	}
	csvToUse, err := getCSV(csvForCalorieSources)
	if err != nil {
		return
	}

	csvFileData, err := getCSVFileData("", csvToUse)
	if err != nil {
		return
	}

	gottenCalorieSourcesData, err := getCalorieSourcesData(calorieSources, csvFileData)
	if err != nil {
		return
	}

	gottenCalorieSourcesData.print()

}

// homemade(homemade);ordered(eatfit,ordered,khannas,dominos);milk-and-cereal(milk)
func validateCalorieSourcesTags(calorieSources string) (map[string][]string, bool) {
	notEnoughTagsErr := "at least two calorie sources tags must be provided (sample format: homemade(homemade);ordered(eatfit,ordered,khannas,dominos);milk-and-cereal(milk))"

	calorieSources = strings.TrimSpace(calorieSources)
	sourceTags := strings.Split(calorieSources, ";")
	if len(sourceTags) < 2 {
		log.Error(notEnoughTagsErr)
		return nil, false
	}

	calorieSourcesMap := make(map[string][]string, 0)

	for _, sourceTag := range sourceTags {

		tagName, subTags, err := validateSourceTag(sourceTag)
		if err != nil {
			return nil, false
		}

		if _, ok := calorieSourcesMap[tagName]; ok {
			log.WithFields(log.Fields{"err": errors.New("a tag name was specified more than once").Error(), "tag_name": tagName}).Error("validation of source tag failed")
			return nil, false
		}

		calorieSourcesMap[tagName] = subTags

	}
	return calorieSourcesMap, true
}

func validateSourceTag(sourceTag string) (string, []string, error) {

	sourceTagValidationErr := "validation of source tag failed"

	if strings.Count(sourceTag, "(") != 1 || strings.Count(sourceTag, ")") != 1 {
		err := errors.New("exactly one starting bracket and one closing bracket were not found")
		log.WithFields(log.Fields{"err": err.Error(), "source_tag": sourceTag}).Error(sourceTagValidationErr)
		return "", nil, err
	}

	subTagsIndex := strings.IndexByte(sourceTag, '(')
	if subTagsIndex == -1 {
		err := errors.New("the first starting bracket was not found as expected")
		log.WithFields(log.Fields{"err": err.Error(), "source_tag": sourceTag}).Error(sourceTagValidationErr)
		return "", nil, err
	}
	if sourceTag[len(sourceTag)-1] != ')' {
		err := errors.New("the source tag does not end with a closing bracket")
		log.WithFields(log.Fields{"err": err.Error(), "source_tag": sourceTag}).Error(sourceTagValidationErr)
		return "", nil, err
	}

	tagName := sourceTag[:subTagsIndex]
	if tagName == "" {
		err := errors.New("the tag name cannot be empty")
		log.WithFields(log.Fields{"err": err.Error(), "source_tag": sourceTag}).Error(sourceTagValidationErr)
		return "", nil, err
	}
	subTags := strings.Split(sourceTag[subTagsIndex+1:len(sourceTag)-1], ",")

	if len(subTags) == 0 {
		err := errors.New("at least one sub tag must be provided")
		log.WithFields(log.Fields{"err": err.Error(), "source_tag": sourceTag}).Error(sourceTagValidationErr)
		return "", nil, err
	}

	return tagName, subTags, nil

}

func getCalorieSourcesData(calorieSources map[string][]string, fileData []map[string]string) (*calorieSourceDataList, error) {

	calorieSourcesData := []calorieSourceData{}
	totalCalories, err := getTotalCalories(fileData)
	if err != nil {
		return nil, err
	}

	for tagName, subTags := range calorieSources {

		calorieSourceDataPoint, err := getDataForSingleCalorieSource(subTags, fileData, totalCalories)
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error(), "tag_name": tagName, "sub_tags": subTags}).Error("failed to get calorie data corresponding to a particular tag")
			return nil, err
		}
		calorieSourcesData = append(calorieSourcesData, calorieSourceData{source: tagName, data: calorieSourceDataPoint})
	}

	return &calorieSourceDataList{data: calorieSourcesData}, nil

}

func getDataForSingleCalorieSource(tags []string, fileData []map[string]string, totalCalories int) (int, error) {

	sumOfTagCalories := 0

	for _, dataPoint := range fileData {
		done := map[string]bool{}
		for _, tag := range tags {

			dataPointStr := strings.ToLower(dataPoint[headerPrefix+"1"])

			if done[dataPointStr] {
				continue
			}

			if strings.Contains(dataPointStr, tag) {

				calories := dataPoint[headerPrefix+"2"]

				caloriesInt, err := strconv.Atoi(calories)
				if err != nil {
					log.WithFields(log.Fields{"err": err.Error(), "data_point": dataPoint}).Error("failed to convert calories from CSV to number")
					return 0, err
				}
				sumOfTagCalories += caloriesInt
				done[dataPointStr] = true

			}
		}
	}

	return 100 * sumOfTagCalories / totalCalories, nil
}

func getTotalCalories(fileData []map[string]string) (int, error) {

	sum := 0

	for _, dataPoint := range fileData {

		calories := dataPoint[headerPrefix+"2"]

		caloriesInt, err := strconv.Atoi(calories)
		if err != nil {
			log.WithFields(log.Fields{"err": err.Error(), "data_point": dataPoint}).Error("failed to convert calories from CSV to number")
			return 0, err
		}

		sum += caloriesInt

	}

	return sum, nil
}

func (c *calorieSourceDataList) print() {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Calorie Source", "Percentage of Calories"})

	for _, calorieData := range c.data {
		row := []string{calorieData.source, strconv.Itoa(calorieData.data)}
		table.Append(row)
	}
	table.Render()
}
func setupCalorieSourcesCmd() {
	csvForCalorieSources = calorieSourcesCmd.Flags().StringP("csv", "f", "", "the path to the CSV to use (optional if a CSV has already been added)")
	calorieSourceTags = calorieSourcesCmd.Flags().StringP("tags", "t", "", `the tags to indicate calorie sources (these tags must be part of the names of food items in the calorie counter CSV); eg. "homemade(homemade);ordered(eatfit,ordered,khannas,dominos);milk-and-cereal(milk)"`)

}
