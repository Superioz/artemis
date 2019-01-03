package config

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	dirUnix    = "/etc/artemis"
	dirWindows = "%s/.artemis"

	configFile = "config.json"

	windowsOS      = "windows"
	homeDriveEnv   = "HOMEDRIVE"
	homePathEnv    = "HOMEPATH"
	userProfileEnv = "USERPROFILE"

	defaultLogDisplayTime = true
	defaultLogDebug       = true

	defaultHost              = "amqp://guest:guest@localhost"
	defaultPort              = "5672"
	defaultExchange          = "artemis"
	defaultBroadcastRoute    = "broadcast.all"
	defaultHeartbeatInterval = 1500
	defaultElectionTimeout   = 2000
	defaultClusterSize       = 2
)

type broker struct {
	Host           string `json:"host"`
	Port           string `json:"port"`
	ExchangeKey    string `json:"exchangeKey"`
	BroadcastRoute string `json:"broadcastRoute"`
}

type NodeConfig struct {
	Broker            broker `json:"broker"`
	Logging           config `json:"logging"`
	HeartbeatInterval int    `json:"heartbeatInterval"`
	ElectionTimeout   int    `json:"electionTimeout"`
	ClusterSize       int    `json:"clusterSize"`
}

func GetRootDirectory() string {
	if runtime.GOOS == windowsOS {
		home := os.Getenv(homeDriveEnv) + os.Getenv(homePathEnv)
		if home == "" {
			home = os.Getenv(userProfileEnv)
		}
		return fmt.Sprintf(dirWindows, strings.Replace(home, "\\", "/", -1))
	}
	return dirUnix
}

func Load() (NodeConfig, error) {
	def := NodeConfig{
		Broker: broker{
			Host:           defaultHost,
			Port:           defaultPort,
			ExchangeKey:    defaultExchange,
			BroadcastRoute: defaultBroadcastRoute,
		},
		Logging: config{
			DisplayTimeStamp: defaultLogDisplayTime,
			Debug:            defaultLogDebug,
		},
		HeartbeatInterval: defaultHeartbeatInterval,
		ElectionTimeout:   defaultElectionTimeout,
		ClusterSize:       defaultClusterSize,
	}

	dir := GetRootDirectory() + "/" + configFile
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(GetRootDirectory(), os.ModePerm)
		if err != nil {
			return def, err
		}

		_, err = os.Create(dir)
		if err != nil {
			return def, err
		}
	}

	c, err := readInConfig(dir, def)
	if err != nil {
		return def, err
	}

	// apply config
	applyConfig(def.Logging)
	return c, nil
}

// reads config struct from given file directory
func readInConfig(file string, config NodeConfig) (NodeConfig, error) {
	configFile, err := os.Open(file)
	if err != nil {
		return config, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		return config, err
	}

	err = configFile.Close()
	return config, err
}
