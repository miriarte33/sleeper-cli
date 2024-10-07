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
	teamOneMap := make(map[int64]MatchupTeam)
	for _, matchupTeamDto := range matchupTeamDtos {
		matchupId := matchupTeamDto.MatchupID
		teamOne, foundTeamOne := teamOneMap[matchupId]

		if foundTeamOne {
			// If we already found teamOne, add both teams to the matchupMap
			teamTwo := mapToMatchupTeam(matchupTeamDto, users, rosters, playerDtos)
			matchupMap[matchupId] = Matchup{
				TeamOne: teamOne,
				TeamTwo: teamTwo,
			}
		} else {
			// Otherwise, add it to the teamOneMap
			teamOneMap[matchupId] = mapToMatchupTeam(matchupTeamDto, users, rosters, playerDtos)
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
	var starterMap = make(map[string]MatchupPlayer)

	var starters []MatchupPlayer
	var bench []MatchupPlayer

	// Populate all the players in the matchup to a Map
	var matchupPlayerMap = make(map[string]MatchupPlayer)
	for playerID, playerPoints := range matchupTeamDto.PlayersPoints {
		player := findPlayer(playerID, playerPoints, playerDtos)
		matchupPlayerMap[playerID] = player
	}

	// Populate the starters in the matchup to a separate Map
	for _, starterID := range matchupTeamDto.Starters {
		player := matchupPlayerMap[starterID]
		starterMap[starterID] = player
		starters = append(starters, player)
	}

	// The bench players are the players in the matchup that are NOT in the starter map
	for _, playerID := range matchupTeamDto.Players {
		_, isStarter := starterMap[playerID]
		if !isStarter {
			player := matchupPlayerMap[playerID]
			bench = append(bench, player)
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
	var rosterOwnerID string
	for _, roster := range rosters {
		if roster.ID != rosterID {
			continue
		}

		rosterOwnerID = roster.OwnerID
	}

	for _, user := range users {
		if user.UserID == rosterOwnerID {
			return &user
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
