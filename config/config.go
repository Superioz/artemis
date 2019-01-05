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
	defaultPort              = 5672
	defaultExchange          = "artemis"
	defaultBroadcastRoute    = "broadcast.all"
	defaultHeartbeatInterval = 1500
	defaultElectionTimeout   = 2000
	defaultClusterSize       = 2

	defaultRestHost            = "localhost"
	defaultRestMinPort         = 2310
	defaultRestMaxPort         = 2315
	defaultRestReadTimeout     = 5
	defaultRestWriteTimeout    = 5
	defaultRestMaxConnsPerIP   = 500
	defaultRestMaxKeepaliveDur = 5
)

// config for the broker connection.
type Broker struct {
	Host           string `json:"host"`
	Port           int    `json:"port"`
	ExchangeKey    string `json:"exchangeKey"`
	BroadcastRoute string `json:"broadcastRoute"`
}

// config for a rest server
// either the internal or external rest api.
type Rest struct {
	Host                 string `json:"host"`
	MinPort              int    `json:"minPort"`
	MaxPort              int    `json:"maxPort"`
	ReadTimeout          int    `json:"readTimeout"`
	WriteTimeout         int    `json:"writeTimeout"`
	MaxConnsPerIP        int    `json:"maxConnsPerIP"`
	MaxKeepaliveDuration int    `json:"maxKeepaliveDuration"`
}

// the common node config.
type NodeConfig struct {
	Broker            Broker  `json:"broker"`
	Logging           Logging `json:"logging"`
	Rest              Rest    `json:"rest"`
	HeartbeatInterval int     `json:"heartbeatInterval"`
	ElectionTimeout   int     `json:"electionTimeout"`
	ClusterSize       int     `json:"clusterSize"`
}

// the root directory for configurations
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

func DefaultBrokerConfig() Broker {
	return Broker{
		Host:           defaultHost,
		Port:           defaultPort,
		ExchangeKey:    defaultExchange,
		BroadcastRoute: defaultBroadcastRoute,
	}
}

func DefaultLoggingConfig() Logging {
	return Logging{
		DisplayTimeStamp: defaultLogDisplayTime,
		Debug:            defaultLogDebug,
	}
}

func DefaultRestConfig() Rest {
	return Rest{
		Host:                 defaultRestHost,
		MinPort:              defaultRestMinPort,
		MaxPort:              defaultRestMaxPort,
		ReadTimeout:          defaultRestReadTimeout,
		WriteTimeout:         defaultRestWriteTimeout,
		MaxConnsPerIP:        defaultRestMaxConnsPerIP,
		MaxKeepaliveDuration: defaultRestMaxKeepaliveDur,
	}
}

// loads the default config or the config from the
// configuration path. returns an error, if something
// went wrong.
func Load() (NodeConfig, error) {
	def := NodeConfig{
		Broker:            DefaultBrokerConfig(),
		Logging:           DefaultLoggingConfig(),
		Rest:              DefaultRestConfig(),
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
