/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miriarte33/sleeper/api"
	envLoader "miriarte33/sleeper/env_loader"
	matchupMapper "miriarte33/sleeper/matchup_mapper"

	"github.com/spf13/cobra"
)

// matchupsCmd represents the matchups command
var matchupsCmd = &cobra.Command{
	Use:   "matchups",
	Short: "View matchups in your league for a given week",
	Long: `View matchups in your league for a given week. Pass in the week with the --week flag.
Note: Matchups that are still in-progress may not have accurate scores because of limitations with Sleeper's free API.`,
	Run: func(cmd *cobra.Command, args []string) {
		week, err := cmd.Flags().GetInt("week")
		if err != nil {
			panic(err)
		}

		if week <= 0 {
			fmt.Println("Week is required. Pass it in with the -w flag.")
			return
		}

		leagueId := envLoader.GetLeagueId()

		// Get matchups
		matchups, err := api.GetMatchupTeams(leagueId, week)
		if err != nil {
			panic(err)
		}

		// Get rosters
		rosters, err := api.GetRosters(leagueId)
		if err != nil {
			panic(err)
		}

		// Get users
		users, err := api.GetUsers(leagueId)
		if err != nil {
			panic(err)
		}

		matchupDict := matchupMapper.MapToMatchupDict(matchups, users, rosters)

		fmt.Println("Matchups for week", week)
		for _, matchup := range matchupDict {
			fmt.Println(
				getMatchupTeamInfo(matchup.TeamOne),
				" vs ",
				getMatchupTeamInfo(matchup.TeamTwo),
			)
		}
	},
}

func getMatchupTeamInfo(team matchupMapper.MatchupTeam) string {
	return fmt.Sprintf("%s (%s) %.2f", team.TeamName, team.UserName, team.TotalPoints)
}

func init() {
	rootCmd.AddCommand(matchupsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// matchupsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// matchupsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	matchupsCmd.Flags().IntP(
		"week",
		"w",
		-1,
		"Required. The week to display matchups for.",
	)
}
