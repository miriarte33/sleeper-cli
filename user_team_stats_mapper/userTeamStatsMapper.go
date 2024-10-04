package userTeamStatsMapper

import (
	"miriarte33/sleeper/api"
	"regexp"
)

type UserTeamStats struct {
	UserID      string
	UserName    string
	TeamName    string
	Wins        int
	Losses      int
	Ties        int
	Fpts        float64
	FptsAgainst float64
}

func MapToUserTeamStatsList(rosters []api.RosterDto, users []api.UserDto) []UserTeamStats {
	var userTeamStats []UserTeamStats
	for _, roster := range rosters {
		for _, user := range users {
			if roster.OwnerID == user.UserID {
				userTeamStats = append(userTeamStats, UserTeamStats{
					UserID:      user.UserID,
					UserName:    user.DisplayName,
					TeamName:    removeEmojis(user.Metadata.TeamName),
					Wins:        roster.Settings.Wins,
					Losses:      roster.Settings.Losses,
					Ties:        roster.Settings.Ties,
					Fpts:        roster.Settings.Fpts,
					FptsAgainst: roster.Settings.FptsAgainst,
				})
				break
			}
		}
	}
	return userTeamStats
}

func removeEmojis(input string) string {
	// Regular expression pattern to match emojis and certain special characters
	emojiRegex := regexp.MustCompile(`[\x{1F600}-\x{1F64F}]|[\x{1F300}-\x{1F5FF}]|[\x{1F680}-\x{1F6FF}]|[\x{2600}-\x{26FF}]|[\x{2700}-\x{27BF}]`)

	// Replace all matched emojis with an empty string
	return emojiRegex.ReplaceAllString(input, "")
}
