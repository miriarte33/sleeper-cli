/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miriarte33/sleeper/api"
	envLoader "miriarte33/sleeper/env_loader"
	matchupMapper "miriarte33/sleeper/matchup_mapper"
	"os"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
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

		preselectedUser, err := cmd.Flags().GetString("username")
		if err != nil {
			panic(err)
		}

		selectedMatchupId := int64(-1)
		if preselectedUser == "" {
			// Map matchup info to a dictionary where the key is the matchup info string and the value is the matchup ID
			matchupInfoOptionsDict := getMatchupInfoOptionsDict(matchupDict)

			// Convert the dictionary keys to a list of keys for the survey prompt
			var matchupInfoOptions []string
			for matchupInfo := range matchupInfoOptionsDict {
				matchupInfoOptions = append(matchupInfoOptions, matchupInfo)
			}

			var selectedMatchupOption string
			prompt := &survey.Select{
				Message: fmt.Sprintf("Matchups for week %d. Select a matchup for details.", week),
				Options: matchupInfoOptions,
			}
			survey.AskOne(prompt, &selectedMatchupOption)

			// Get the matchup ID from the selectedMatchupOption by key
			selectedMatchupId = matchupInfoOptionsDict[selectedMatchupOption]
		} else {
			// Find the matchup by the preselected user
			for matchupId, matchup := range matchupDict {
				if matchup.TeamOne.UserName == preselectedUser || matchup.TeamTwo.UserName == preselectedUser {
					selectedMatchupId = matchupId
					break
				}
			}
		}

		// Print the selected matchup
		if selectedMatchupId == -1 {
			fmt.Println("No matchup found for the given week and username.")
			return
		}

		matchup := matchupDict[selectedMatchupId]

		scoreWriter := tabwriter.NewWriter(os.Stdout, 0, 5, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
		fmt.Fprintf(
			scoreWriter,
			"%s (%s)\t%.2f\t%.2f\t%s (%s)\t\n",
			matchup.TeamOne.TeamName,
			matchup.TeamOne.UserName,
			matchup.TeamOne.TotalPoints,
			matchup.TeamTwo.TotalPoints,
			matchup.TeamTwo.TeamName,
			matchup.TeamTwo.UserName,
		)

		scoreWriter.Flush()

		fmt.Println()
		startersWriter := tabwriter.NewWriter(os.Stdout, 0, 5, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
		fmt.Fprintf(
			startersWriter,
			"Starters (%s)\t-\t-\tStarters (%s)\t\n",
			matchup.TeamOne.UserName,
			matchup.TeamTwo.UserName,
		)

		formattedStarterComparisons := formatPlayersSideBySide(
			matchup.TeamOne.Starters,
			matchup.TeamTwo.Starters,
		)
		for _, formattedPlayerComparison := range formattedStarterComparisons {
			fmt.Fprint(
				startersWriter,
				formattedPlayerComparison,
			)
		}
		startersWriter.Flush()

		fmt.Println()
		benchWriter := tabwriter.NewWriter(os.Stdout, 0, 5, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)
		fmt.Fprintf(
			benchWriter,
			"Bench (%s)\t-\t-\tBench (%s)\t\n",
			matchup.TeamOne.UserName,
			matchup.TeamTwo.UserName,
		)

		formattedBenchComparisons := formatPlayersSideBySide(
			matchup.TeamOne.Bench,
			matchup.TeamTwo.Bench,
		)
		for _, formattedBenchComparison := range formattedBenchComparisons {
			fmt.Fprint(
				benchWriter,
				formattedBenchComparison,
			)
		}
		benchWriter.Flush()
	},
}

// Formats the players side by side to compare their points
func formatPlayersSideBySide(list1, list2 []matchupMapper.MatchupPlayer) []string {
	maxLength := len(list1)
	if len(list2) > maxLength {
		maxLength = len(list2)
	}

	var formattedPlayers []string
	for i := 0; i < maxLength; i++ {
		var str1, str2 string

		if i < len(list1) {
			player1 := list1[i]
			str1 = fmt.Sprintf(
				"%s (%s)\t%.2f",
				player1.FullName,
				player1.Position,
				player1.Points,
			)
		} else {
			str1 = "-\t-\t"
		}

		if i < len(list2) {
			player1 := list2[i]
			str2 = fmt.Sprintf(
				"%.2f\t%s (%s)",
				player1.Points,
				player1.FullName,
				player1.Position,
			)
		} else {
			str2 = "-\t-"
		}

		formattedPlayers = append(
			formattedPlayers,
			fmt.Sprintf("%s\t%s\t\n", str1, str2),
		)
	}

	return formattedPlayers
}

// Returns a map where the Key is the matchup info string and the value is the matchup ID
func getMatchupInfoOptionsDict(matchupDict map[int64]matchupMapper.Matchup) map[string]int64 {
	var teamInfoList = make(map[string]int64)
	for matchupId, matchup := range matchupDict {
		teamInfoList[getMatchupResults(matchup)] = matchupId
	}
	return teamInfoList
}

func getMatchupResults(matchup matchupMapper.Matchup) string {
	return fmt.Sprint(
		getMatchupTeamInfo(matchup.TeamOne),
		" vs ",
		getMatchupTeamInfo(matchup.TeamTwo),
	)
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

	matchupsCmd.Flags().StringP(
		"username",
		"u",
		"",
		"Optionally specify the username whose matchup to view for the given week.",
	)
}
