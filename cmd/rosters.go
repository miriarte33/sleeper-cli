/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miriarte33/sleeper/api"
	envLoader "miriarte33/sleeper/env_loader"
	playerLoader "miriarte33/sleeper/player_loader"
	"os"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

// rostersCmd represents the rosters command
var rostersCmd = &cobra.Command{
	Use:   "rosters",
	Short: "View rosters in your league",
	Long:  `View rosters in your league`,
	Run: func(cmd *cobra.Command, args []string) {
		username, err := cmd.Flags().GetString("username")
		if err != nil {
			panic(err)
		}

		leagueId := envLoader.GetLeagueId()
		users, err := api.GetUsers(leagueId)

		if err != nil {
			panic(err)
		}

		options := getUserNameList(users)

		var selectedUserName string
		if username != "" {
			selectedUserName = username
		} else {
			prompt := &survey.Select{
				Message: "Select a user's roster to view:",
				Options: options,
			}
			survey.AskOne(prompt, &selectedUserName)
		}

		selectedUser := findUserByDisplayName(users, selectedUserName)
		if selectedUser == nil {
			fmt.Printf("User not found in your league: %s\n", selectedUserName)
			return
		}

		rosters, err := api.GetRosters(leagueId)
		if err != nil {
			panic(err)
		}

		selectedUserRoster := findRosterByUserId(rosters, selectedUser.UserID)

		printTeamStatsForRoster(selectedUserRoster)
		fmt.Println()
		printPlayersInRoster(selectedUserRoster)
	},
}

func printPlayersInRoster(selectedUserRoster *api.RosterDto) {
	playersInRoster, err := getPlayersInRoster(*selectedUserRoster)
	if err != nil {
		panic(err)
	}

	// Create a new tab writer
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)

	// Table columns
	fmt.Fprintln(writer, "Name\tPosition\tTeam\tAge\tInjury Status\tInjury\t")

	// Write the table rows
	for _, player := range playersInRoster {
		if player.Position == "QB" {
			printPlayerInfo(writer, player)
		}
	}

	for _, player := range playersInRoster {
		if player.Position == "RB" {
			printPlayerInfo(writer, player)
		}
	}

	for _, player := range playersInRoster {
		if player.Position == "WR" {
			printPlayerInfo(writer, player)
		}
	}

	for _, player := range playersInRoster {
		if player.Position == "TE" {
			printPlayerInfo(writer, player)
		}
	}

	// Flush the writer to output the table
	writer.Flush()
}

func printPlayerInfo(writer *tabwriter.Writer, player playerLoader.PlayerDto) {
	fantasyPositions := ""
	for i, pos := range player.FantasyPositions {
		if i > 0 {
			fantasyPositions += ", "
		}
		fantasyPositions += pos
	}

	fmt.Fprintf(
		writer,
		"%s\t%s\t%s\t%d\t%s\t%s\t\n",
		player.FullName,
		fantasyPositions,
		player.Team,
		player.Age,
		player.InjuryStatus,
		player.InjuryBodyPart,
	)
}

func getPlayersInRoster(rosterDto api.RosterDto) ([]playerLoader.PlayerDto, error) {
	players, err := playerLoader.LoadPlayers()
	if err != nil {
		return nil, fmt.Errorf("failed to load players: %w", err)
	}

	var rosterPlayers []playerLoader.PlayerDto
	for _, playerID := range rosterDto.Players {
		player, ok := players[playerID]
		if !ok {
			return nil, fmt.Errorf("player not found: %s", playerID)
		}
		rosterPlayers = append(rosterPlayers, player)
	}

	return rosterPlayers, nil
}

func printTeamStatsForRoster(selectedUserRoster *api.RosterDto) {
	// Create a new tab writer
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)

	// Table columns
	fmt.Fprintln(writer, "Wins\tLosses\tTies\tPts For\tPts Against\t")

	// Write the table rows
	fmt.Fprintf(
		writer,
		"%d\t%d\t%d\t%.2f\t%.2f\t\n",
		selectedUserRoster.Settings.Wins,
		selectedUserRoster.Settings.Losses,
		selectedUserRoster.Settings.Ties,
		selectedUserRoster.Settings.Fpts,
		selectedUserRoster.Settings.FptsAgainst,
	)

	// Flush the writer to output the table
	writer.Flush()
}

func findUserByDisplayName(users []api.UserDto, displayName string) *api.UserDto {
	for _, user := range users {
		if user.DisplayName == displayName {
			return &user
		}
	}
	return nil
}

func findRosterByUserId(rosters []api.RosterDto, userId string) *api.RosterDto {
	for _, roster := range rosters {
		if roster.OwnerID == userId {
			return &roster
		}
	}
	return nil
}

func getUserNameList(users []api.UserDto) []string {
	userNames := make([]string, len(users))
	for i, user := range users {
		userNames[i] = user.DisplayName
	}
	return userNames
}

func init() {
	rootCmd.AddCommand(rostersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rostersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	rostersCmd.Flags().StringP(
		"username",
		"u",
		"",
		"Optionally specify the username whose roster to view",
	)
}
