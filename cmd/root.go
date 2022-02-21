package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "calorie-insights",
	Short: "get insights from a calorie CSV",
	Long:  `get insights from a calorie CSV`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("error executing root command")
		os.Exit(1)
	}
}

func init() {

	setupWeeklyAvgCmd()
	setupAddCSVCmd()
	rootCmd.AddCommand(weeklyAvgCmd)
	rootCmd.AddCommand(addCSVCmd)

}
