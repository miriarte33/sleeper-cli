package playerLoader

import (
	"encoding/json"
	"fmt"
	"miriarte33/sleeper/api"
	"os"
)

type PlayerDto struct {
	FullName         string   `json:"full_name"`
	Position         string   `json:"position"`
	Team             string   `json:"team"`
	Age              int      `json:"age"`
	InjuryBodyPart   string   `json:"injury_body_part"`
	FantasyPositions []string `json:"fantasy_positions"`
	InjuryStatus     string   `json:"injury_status"`
	Status           string   `json:"status"`
}

func GetPlayersInRoster(rosterDto api.RosterDto) ([]PlayerDto, error) {
	players, err := LoadPlayers()
	if err != nil {
		return nil, fmt.Errorf("failed to load players: %w", err)
	}

	var rosterPlayers []PlayerDto
	for _, playerID := range rosterDto.Players {
		player, ok := players[playerID]
		if !ok {
			return nil, fmt.Errorf("player not found: %s", playerID)
		}
		rosterPlayers = append(rosterPlayers, player)
	}

	return rosterPlayers, nil
}

func LoadPlayers() (map[string]PlayerDto, error) {
	data, err := os.ReadFile("players.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	players := make(map[string]PlayerDto)
	err = json.Unmarshal(data, &players)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return players, nil
}
