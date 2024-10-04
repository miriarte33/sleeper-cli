/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"miriarte33/sleeper/api"
	envLoader "miriarte33/sleeper/env_loader"
	userTeamStatsMapper "miriarte33/sleeper/user_team_stats_mapper"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

// standingsCmd represents the standings command
var standingsCmd = &cobra.Command{
	Use:   "standings",
	Short: "View your leagues standings",
	Long:  `View your leagues standings`,
	Run: func(cmd *cobra.Command, args []string) {
		leagueId := envLoader.GetLeagueId()
		users, err := api.GetUsers(leagueId)
		if err != nil {
			panic(err)
		}

		rosters, err := api.GetRosters(leagueId)
		if err != nil {
			panic(err)
		}

		userTeamStats := userTeamStatsMapper.MapToUserTeamStatsList(rosters, users)

		sort.Slice(userTeamStats, func(i, j int) bool {
			if userTeamStats[i].Wins == userTeamStats[j].Wins {
				return userTeamStats[i].Fpts > userTeamStats[j].Fpts
			}
			return userTeamStats[i].Wins > userTeamStats[j].Wins
		})

		// Create a new tab writer
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug|tabwriter.AlignRight)

		// Table columns
		fmt.Fprintln(writer, "Owner\tTeam Name\tWins\tLosses\tTies\tPts For\tPts Against\t")

		for _, userTeamStats := range userTeamStats {
			// Write the table rows
			fmt.Fprintf(
				writer,
				"%s\t%s\t%d\t%d\t%d\t%.2f\t%.2f\t\n",
				userTeamStats.UserName,
				userTeamStats.TeamName,
				userTeamStats.Wins,
				userTeamStats.Losses,
				userTeamStats.Ties,
				userTeamStats.Fpts,
				userTeamStats.FptsAgainst,
			)
		}

		writer.Flush()
	},
}

func init() {
	rootCmd.AddCommand(standingsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// standingsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// standingsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
