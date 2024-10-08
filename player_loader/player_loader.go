package playerLoader

import (
	"encoding/json"
	"fmt"
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
