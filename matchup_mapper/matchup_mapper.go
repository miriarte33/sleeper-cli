package matchupMapper

import (
	"fmt"
	"miriarte33/sleeper/api"
	playerLoader "miriarte33/sleeper/player_loader"
)

type MatchupPlayer struct {
	PlayerID string
	Points   float64
	FullName string
	Position string
}

type MatchupTeam struct {
	UserID      string
	UserName    string
	TeamName    string
	TotalPoints float64
	Starters    []MatchupPlayer
	Bench       []MatchupPlayer
}

type Matchup struct {
	TeamOne MatchupTeam
	TeamTwo MatchupTeam
}

func MapToMatchupDict(
	matchupTeamDtos []api.MatchupTeamDto,
	users []api.UserDto,
	rosters []api.RosterDto,
) map[int64]Matchup {
	playerDtos, err := playerLoader.LoadPlayers()

	if err != nil {
		panic(err)
	}

	matchupMap := make(map[int64]Matchup)

	for _, matchupTeamDto := range matchupTeamDtos {
		matchupId := matchupTeamDto.MatchupID
		_, exists := matchupMap[matchupId]

		// If the matchup doesn't exist, create it
		if exists {
			continue
		}

		teamOne := mapToMatchupTeam(matchupTeamDto, users, rosters, playerDtos)

		// Loop through the matchupTeams to find the second team
		for _, matchupTeam2Dto := range matchupTeamDtos {
			if matchupTeam2Dto.MatchupID != matchupId {
				continue
			}

			teamTwo := mapToMatchupTeam(matchupTeam2Dto, users, rosters, playerDtos)

			matchupMap[matchupId] = Matchup{
				TeamOne: teamOne,
				TeamTwo: teamTwo,
			}
		}
	}

	return matchupMap
}

func mapToMatchupTeam(
	matchupTeamDto api.MatchupTeamDto,
	users []api.UserDto,
	rosters []api.RosterDto,
	playerDtos map[string]playerLoader.PlayerDto,
) MatchupTeam {
	var starters []MatchupPlayer
	var bench []MatchupPlayer

	for playerID, playerPoints := range matchupTeamDto.PlayersPoints {
		player := findPlayer(playerID, playerPoints, playerDtos)

		// If the player is NOT in matchupTeamDto.Starters, add them to the bench
		if !isStarter(playerID, matchupTeamDto.Starters) {
			bench = append(bench, player)
		} else {
			starters = append(starters, player)
		}
	}

	user := findUser(matchupTeamDto.RosterID, users, rosters)

	if user == nil {
		panic("user not found: " + fmt.Sprint(matchupTeamDto.MatchupID))
	}

	return MatchupTeam{
		UserID:      user.UserID,
		UserName:    user.DisplayName,
		TeamName:    user.Metadata.TeamName,
		TotalPoints: matchupTeamDto.Points,
		Starters:    starters,
		Bench:       bench,
	}
}

func findUser(
	rosterID int64,
	users []api.UserDto,
	rosters []api.RosterDto,
) *api.UserDto {
	for _, roster := range rosters {
		if roster.ID != rosterID {
			continue
		}

		for _, user := range users {
			if roster.OwnerID == user.UserID {
				return &user
			}
		}
	}

	return nil
}

func findPlayer(
	playerID string,
	playerPoints float64,
	players map[string]playerLoader.PlayerDto,
) MatchupPlayer {
	player, ok := players[playerID]
	if !ok {
		panic("player not found: " + playerID)
	}

	return MatchupPlayer{
		PlayerID: playerID,
		Points:   playerPoints,
		FullName: player.FullName,
		Position: player.Position,
	}
}

func isStarter(
	playerID string,
	starters []string,
) bool {
	for _, starterID := range starters {
		if playerID == starterID {
			return true
		}
	}
	return false
}
