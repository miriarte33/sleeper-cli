package envLoader

import (
	"os"

	"gopkg.in/yaml.v2"
)

func GetLeagueId() string {
	err := loadEnvFromFile("env.yaml")
	if err != nil {
		panic(err)
	}

	leagueId := os.Getenv("LEAGUE_ID")
	if leagueId == "" {
		panic("LEAGUE_ID environment variable not set")
	}

	return leagueId
}

func loadEnvFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	env := make(map[string]string)
	err = yaml.Unmarshal(data, &env)
	if err != nil {
		return err
	}

	for key, value := range env {
		os.Setenv(key, value)
	}

	return nil
}
