/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miriarte33/sleeper/api"
	envLoader "miriarte33/sleeper/env_loader"
	matchupMapper "miriarte33/sleeper/matchup_mapper"

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

		fmt.Println(matchupDict[selectedMatchupId])
	},
}

// Returns a map where the Key is the matchup info string and the value is the matchup ID
func getMatchupInfoOptionsDict(matchupDict map[int64]matchupMapper.Matchup) map[string]int64 {
	var teamInfoList = make(map[string]int64)
	for matchupId, matchup := range matchupDict {
		teamInfoList[fmt.Sprint(
			getMatchupTeamInfo(matchup.TeamOne),
			" vs ",
			getMatchupTeamInfo(matchup.TeamTwo),
		)] = matchupId
	}
	return teamInfoList
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
