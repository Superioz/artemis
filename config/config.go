package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	defaultRestPort            = 2310
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
type CLIRest struct {
	Host                 string `json:"host"`
	Port                 int    `json:"port"`
	ReadTimeout          int    `json:"readTimeout"`
	WriteTimeout         int    `json:"writeTimeout"`
	MaxConnsPerIP        int    `json:"maxConnsPerIP"`
	MaxKeepaliveDuration int    `json:"maxKeepaliveDuration"`
}

// the common node config.
type NodeConfig struct {
	Broker            Broker  `json:"broker"`
	Logging           Logging `json:"logging"`
	CLIRest           CLIRest `json:"clirest"`
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

func DefaultCLIRestConfig() CLIRest {
	return CLIRest{
		Host:                 defaultRestHost,
		Port:                 defaultRestPort,
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
		CLIRest:           DefaultCLIRestConfig(),
		HeartbeatInterval: defaultHeartbeatInterval,
		ElectionTimeout:   defaultElectionTimeout,
		ClusterSize:       defaultClusterSize,
	}
	c, err := getConfig(def)

	ApplyLoggingConfig(c.Logging)
	return def, err
}

// gets a config object from path
// otherwise return `def`
func getConfig(def NodeConfig) (NodeConfig, error) {
	dir := GetRootDirectory() + "/" + configFile
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if _, err := os.Stat(GetRootDirectory()); os.IsNotExist(err) {
			err := os.Mkdir(GetRootDirectory(), os.ModePerm)
			if err != nil {
				return def, err
			}
		}

		_, err = os.Create(dir)
		if err != nil {
			return def, err
		} else {
			// write default content to prevent `EOF` error
			_ = ioutil.WriteFile(dir, []byte("{}"), os.ModePerm)
		}
	}

	c, err := readInConfig(dir, def)
	if err != nil {
		return def, err
	}

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
