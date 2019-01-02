package config

import (
	"encoding/json"
	"os"
)

type ClusterConfig struct {
	Broker struct {
		Host           string `json:"host"`
		Port           string `json:"port"`
		ExchangeKey    string `json:"exchangeKey"`
		BroadcastRoute string `json:"broadcastRoute"`
	} `json:"broker"`
	HeartbeatInterval int `json:"heartbeatInterval"`
	ElectionTimeout   int `json:"electionTimeout"`
	ClusterSize       int `json:"clusterSize"`
}

func Load(file string) (ClusterConfig, error) {
	var config ClusterConfig
	configFile, err := os.Open(file)
	if err != nil {
		return ClusterConfig{}, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return ClusterConfig{}, err
	}

	err = configFile.Close()
	return config, err
}
