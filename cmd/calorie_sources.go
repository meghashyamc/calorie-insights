package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var calorieSources = &cobra.Command{
	Use:   "sources",
	Short: "get percentages of calories from various sources",
	Long:  "get percentages of calories from various sources",
	Run:   getCalorieSourceDetails,
}

func getCalorieSourceDetails(cmd *cobra.Command, args []string) {

	if len(args) == 0 {
		log.Error("getting calorie source details requires a CSV file path argument that wasn't found | alternatively add a CSV by using the addcsv command")
		return
	}
}
