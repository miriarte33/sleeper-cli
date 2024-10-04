/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miriarte33/sleeper/api"
	envLoader "miriarte33/sleeper/env_loader"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// leagueCmd represents the league command
var leagueCmd = &cobra.Command{
	Use:   "league",
	Short: "Displays your leagues settings",
	Long:  `Display your leagues settings`,
	Run: func(cmd *cobra.Command, args []string) {
		leagueId := envLoader.GetLeagueId()
		league, err := api.GetLeague(leagueId)

		if err != nil {
			panic(err)
		}

		fmt.Printf("League ID: %s\n", league.ID)
		fmt.Printf("Name: %s\n", league.Name)
		fmt.Printf("Season: %s\n", league.Season)
		fmt.Printf("Number of Teams: %d\n", league.TotalRosters)
		bestBallYesNo := "No"
		if league.Settings.BestBall == 1 {
			bestBallYesNo = "Yes"
		}
		fmt.Printf("Best Ball: %s\n", bestBallYesNo)

		fmt.Println()

		fmt.Println("Roster Positions:")
		// Create a new tab writer
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
		fmt.Fprintln(writer, "QB\tRB\tWR\tTE\tFLEX\tBN\t")

		// Write the table rows
		fmt.Fprintf(
			writer,
			"%d\t%d\t%d\t%d\t%d\t%d\t\n",
			countMatches(league.RosterPositions, "QB"),
			countMatches(league.RosterPositions, "RB"),
			countMatches(league.RosterPositions, "WR"),
			countMatches(league.RosterPositions, "TE"),
			countMatches(league.RosterPositions, "FLEX"),
			countMatches(league.RosterPositions, "BN"),
		)

		writer.Flush()

		fmt.Println()

		fmt.Println("Scoring Settings:")
		// Create a new tab writer
		scoringSettingsWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
		fmt.Fprintln(scoringSettingsWriter, "Pass Yd\tPass TD\tPass 2pt\tPass Int\tRush Yd\tRush TD\tRush 2pt\tRec Yd\tRec TD\tRec 2pt\tRec\tFum Lost\t")

		// Write the table rows
		fmt.Fprintf(
			scoringSettingsWriter,
			"%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t%.2f\t\n",
			league.ScoringSettings.PassYd,
			league.ScoringSettings.PassTd,
			league.ScoringSettings.Pass2pt,
			league.ScoringSettings.PassInt,
			league.ScoringSettings.RushYd,
			league.ScoringSettings.RushTd,
			league.ScoringSettings.Rush2pt,
			league.ScoringSettings.RecYd,
			league.ScoringSettings.RecTd,
			league.ScoringSettings.Rec2pt,
			league.ScoringSettings.Rec,
			league.ScoringSettings.FumLost,
		)
		scoringSettingsWriter.Flush()
	},
}

func countMatches(list []string, match string) int {
	count := 0
	for _, item := range list {
		if item == match {
			count++
		}
	}
	return count
}

func init() {
	rootCmd.AddCommand(leagueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// leagueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// leagueCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
