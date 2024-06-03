package config

import (
	"IPFS-CLUSTER-MANAGER/internal/log"
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_Load(t *testing.T) {

	cases := []struct {
		name          string
		setupEnv      func()
		expectedError bool
	}{
		{
			name: "valid configuration",
			setupEnv: func() {
				setEnvironment()
			},
			expectedError: false,
		},
		{
			name: "missing SERVER_URL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("SERVER_URL")
			},
			expectedError: true,
		},
		{
			name: "missing SERVER_PORT",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("SERVER_PORT")
			},
			expectedError: true,
		},
		{
			name: "missing LOG_LEVEL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("LOG_LEVEL")
			},
			expectedError: false,
		},
		{
			name: "wrong value for LOG_LEVEL",
			setupEnv: func() {
				setEnvironment()
				os.Setenv("LOG_LEVEL", "10")
			},
			expectedError: true,
		},
		{
			name: "missing LOG_MODE",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("LOG_MODE")
			},
			expectedError: true,
		},
		{
			name: "wrong value for LOG_MODE",
			setupEnv: func() {
				setEnvironment()
				os.Setenv("LOG_MODE", "3")
			},
			expectedError: true,
		},
		{
			name: "missing IPFS0_NODE_URL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("IPFS0_NODE_URL")
			},
			expectedError: true,
		},
		{
			name: "missing IPFS0_CLUSTER_URL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("IPFS0_CLUSTER_URL")
			},
			expectedError: true,
		},
		{
			name: "missing IPFS1_NODE_URL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("IPFS1_NODE_URL")
			},
			expectedError: true,
		},
		{
			name: "missing IPFS1_CLUSTER_URL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("IPFS1_CLUSTER_URL")
			},
			expectedError: true,
		},
		{
			name: "missing IPFS2_NODE_URL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("IPFS2_NODE_URL")
			},
			expectedError: true,
		},
		{
			name: "missing IPFS2_CLUSTER_URL",
			setupEnv: func() {
				setEnvironment()
				os.Unsetenv("IPFS2_CLUSTER_URL")
			},
			expectedError: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupEnv()
			_, err := Load()
			if tc.expectedError {
				log.Info(context.Background(), err.Error())
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func setEnvironment() {
	os.Setenv("SERVER_URL", "http://localhost")
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("LOG_LEVEL", "-4")
	os.Setenv("LOG_MODE", "2")

	os.Setenv("IPFS0_NODE_URL", "http://localhost:6003")
	os.Setenv("IPFS0_CLUSTER_URL", "http://localhost:7003")
	os.Setenv("IPFS1_NODE_URL", "http://localhost:6001")
	os.Setenv("IPFS1_CLUSTER_URL", "http://localhost:7001")
	os.Setenv("IPFS2_NODE_URL", "http://localhost:6002")
	os.Setenv("IPFS2_CLUSTER_URL", "http://localhost:7002")
}
