package config

import (
	"IPFS-CLUSTER-MANAGER/internal/log"
	"context"
	"errors"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Configuration struct {
	ServerUrl       string `mapstructure:"ServerUrl" tip:"Server Url"`
	ServerPort      int    `mapstructure:"ServerPort" tip:"Server port"`
	Ipfs0NodeUrl    string `mapstructure:"Ipfs0NodeUrl" tip:"Ipfs1 port"`
	Ipfs0ClusterUrl string `mapstructure:"Ipfs0ClusterUrl" tip:"Ipfs1 cluster port"`
	Ipfs1NodeUrl    string `mapstructure:"Ipfs1NodeUrl" tip:"Ipfs2 port"`
	Ipfs1ClusterUrl string `mapstructure:"Ipfs1ClusterUrl" tip:"Ipfs2 cluster port"`
	Ipfs2NodeUrl    string `mapstructure:"Ipfs2NodeUrl" tip:"Ipfs3 port"`
	Ipfs2ClusterUrl string `mapstructure:"Ipfs2ClusterUrl" tip:"Ipfs3 cluster port"`
	Log             Log    `mapstructure:"Log"`
}

type Log struct {
	Level int `mapstructure:"Level" tip:"Minimum level to log: (-4:Debug, 0:Info, 4:Warning, 8:Error)"`
	Mode  int `mapstructure:"Mode" tip:"Log format (1: JSON, 2:Structured text)"`
}

func Load() (*Configuration, error) {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Warn(ctx, "Error loading .env file, using default environment variables", "err", err)
	}

	bindEnv()

	config := &Configuration{}

	if err := viper.Unmarshal(config); err != nil {
		log.Error(ctx, "error unmarshalling configuration", "err", err)
	}

	err = checkEnvVars(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
func bindEnv() {
	_ = viper.BindEnv("ServerUrl", "SERVER_URL")
	_ = viper.BindEnv("ServerPort", "SERVER_PORT")

	_ = viper.BindEnv("Ipfs0NodeUrl", "IPFS0_NODE_URL")
	_ = viper.BindEnv("Ipfs0ClusterUrl", "IPFS0_CLUSTER_URL")
	_ = viper.BindEnv("Ipfs1NodeUrl", "IPFS1_NODE_URL")
	_ = viper.BindEnv("Ipfs1ClusterUrl", "IPFS1_CLUSTER_URL")
	_ = viper.BindEnv("Ipfs2NodeUrl", "IPFS2_NODE_URL")
	_ = viper.BindEnv("Ipfs2ClusterUrl", "IPFS2_CLUSTER_URL")

	_ = viper.BindEnv("Log.Level", "LOG_LEVEL")
	_ = viper.BindEnv("Log.Mode", "LOG_MODE")

	viper.AutomaticEnv()
}

func checkEnvVars(cfg *Configuration) error {
	if cfg.ServerUrl == "" {
		return errors.New("SERVER_URL env var is required")
	}
	if cfg.ServerPort == 0 {
		return errors.New("SERVER_PORT env var is required")
	}
	if cfg.Ipfs0NodeUrl == "" {
		return errors.New("IPFS0_PORT env var is required")
	}
	if cfg.Ipfs0ClusterUrl == "" {
		return errors.New("IPFS0_CLUSTER_PORT env var is required")
	}
	if cfg.Ipfs1NodeUrl == "" {
		return errors.New("IPFS1_PORT env var is required")
	}
	if cfg.Ipfs1ClusterUrl == "" {
		return errors.New("IPFS1_CLUSTER_PORT env var is required")
	}
	if cfg.Ipfs2NodeUrl == "" {
		return errors.New("IPFS2_PORT env var is required")
	}
	if cfg.Ipfs2ClusterUrl == "" {
		return errors.New("IPFS2_CLUSTER_PORT env var is required")
	}
	logLevels := []int{-4, 0, 4, 8}
	if isInSlice(cfg.Log.Level, logLevels) == false {
		return errors.New("LOG_LEVEL env var is required. Possible values: [-4, 0, 4, 8]")
	}
	logModes := []int{1, 2}
	if isInSlice(cfg.Log.Mode, logModes) == false {
		return errors.New("LOG_MODE env var is required. Possible values: [1, 2]")
	}
	return nil
}

func isInSlice(value int, slice []int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
