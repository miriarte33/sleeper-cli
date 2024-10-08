package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LeagueDto struct {
	ID              string   `json:"league_id"`
	Name            string   `json:"name"`
	TotalRosters    int      `json:"total_rosters"`
	Status          string   `json:"status"`
	Sport           string   `json:"sport"`
	SeasonType      string   `json:"season_type"`
	Season          string   `json:"season"`
	RosterPositions []string `json:"roster_positions"`
	Settings        struct {
		BestBall int `json:"best_ball"`
		NumTeams int `json:"num_teams"`
	} `json:"settings"`
	ScoringSettings struct {
		PassInt float64 `json:"pass_int"`
		Pass2pt float64 `json:"pass_2pt"`
		PassYd  float64 `json:"pass_yd"`
		PassTd  float64 `json:"pass_td"`
		Rush2pt float64 `json:"rush_2pt"`
		RushYd  float64 `json:"rush_yd"`
		RushTd  float64 `json:"rush_td"`
		Rec2pt  float64 `json:"rec_2pt"`
		RecYd   float64 `json:"rec_yd"`
		RecTd   float64 `json:"rec_td"`
		Rec     float64 `json:"rec"`
		FumLost float64 `json:"fum_lost"`
	} `json:"scoring_settings"`
}

func GetLeague(leagueID string) (*LeagueDto, error) {
	url := fmt.Sprintf("https://api.sleeper.app/v1/league/%s", leagueID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get league: %s", resp.Status)
	}

	var league LeagueDto
	if err := json.NewDecoder(resp.Body).Decode(&league); err != nil {
		return nil, err
	}

	return &league, nil
}

type RosterDto struct {
	ID       int64    `json:"roster_id"`
	OwnerID  string   `json:"owner_id"`
	Players  []string `json:"players"`
	Starters []string `json:"starters"`
	Settings struct {
		Wins        int     `json:"wins"`
		Losses      int     `json:"losses"`
		Ties        int     `json:"ties"`
		Fpts        float64 `json:"fpts"`
		FptsAgainst float64 `json:"fpts_against"`
	} `json:"settings"`
}

func GetRosters(leagueID string) ([]RosterDto, error) {
	url := fmt.Sprintf("https://api.sleeper.app/v1/league/%s/rosters", leagueID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get rosters: %s", resp.Status)
	}

	var rosters []RosterDto
	if err := json.NewDecoder(resp.Body).Decode(&rosters); err != nil {
		return nil, err
	}

	return rosters, nil
}

type UserDto struct {
	UserID      string `json:"user_id"`
	DisplayName string `json:"display_name"`
	Metadata    struct {
		TeamName string `json:"team_name"`
	} `json:"metadata"`
}

func GetUsers(leagueID string) ([]UserDto, error) {
	url := fmt.Sprintf("https://api.sleeper.app/v1/league/%s/users", leagueID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get users: %s", resp.Status)
	}

	var users []UserDto
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, err
	}

	return users, nil
}

type MatchupTeamDto struct {
	MatchupID     int64              `json:"matchup_id"`
	RosterID      int64              `json:"roster_id"`
	Points        float64            `json:"points"`
	Players       []string           `json:"players"`
	Starters      []string           `json:"starters"`
	PlayersPoints map[string]float64 `json:"players_points"`
}

func GetMatchupTeams(leagueID string, week int) ([]MatchupTeamDto, error) {
	url := fmt.Sprintf("https://api.sleeper.app/v1/league/%s/matchups/%d", leagueID, week)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get matchups: %s", resp.Status)
	}

	var matchups []MatchupTeamDto
	if err := json.NewDecoder(resp.Body).Decode(&matchups); err != nil {
		return nil, err
	}

	return matchups, nil
}
