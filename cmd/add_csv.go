package cmd

import (
	"io/fs"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var csvForAddCSV *string
var addCSVCmd = &cobra.Command{
	Use:   "addcsv",
	Short: "add a calories CSV for analysis",
	Long:  "add a calories CSV for analysis",
	Run:   addCalorieCSV,
}

func addCalorieCSV(cmd *cobra.Command, args []string) {

	if *csvForAddCSV == "" {
		log.Error("a CSV path must be specfied to add CSV")
		return
	}

	fileBytes, err := ioutil.ReadFile(*csvForAddCSV)
	if err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not read from passed filepath")
		return
	}

	if err := ioutil.WriteFile("calorie_counter.csv", fileBytes, fs.ModePerm); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("could not add CSV")
		return
	}

	log.Info("added CSV successfully")

}

func setupAddCSVCmd() {
	csvForAddCSV = addCSVCmd.Flags().StringP("csv", "f", "", "the path to the CSV to be added")
}
