package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	log "github.com/Sirupsen/logrus"
)

func TestGetAPIVersion(t *testing.T) {
	assert.Equal(t, "1",getAPIVersion("1.0.0"), "getting version")
	assert.Equal(t, "0",getAPIVersion("1.0"), "getting version")
}

func TestBuildCommands(t *testing.T) {
	var configurations = []Configuration{
		{
			Url:     "/v1/command1",
			Command: "/opt/command1",
		},
		{
			Url:     "/v1/command2",
			Command: "/opt/command2 -arg1 arg2",
		},
	}

	var commands = map[string]string{
		"/v1/command1": "/opt/command1",
		"/v1/command2": "/opt/command2 -arg1 arg2",
	}
	assert.Equal(t, commands, buildCommands(configurations), "build commands")
}

func TestSetUp(t *testing.T) {

	setUp()
	assert.Equal(t, ":8891",Port,"default port")
	assert.Equal(t, "./configuration.json", ConfigurationFile, "default configuration file")
	assert.Equal(t, "./credentials.json", CredentialsFile, "default credentials")


	os.Setenv("PORT", "123")
	os.Setenv("FILE_CONFIGURATION", "/etc/rest2command/configuration.json")
	os.Setenv("FILE_CREDENTIALS", "/etc/rest2command/credentials.json")
	setUp()
	assert.Equal(t, ":123",Port,"Setting port")
	assert.Equal(t, "/etc/rest2command/configuration.json", ConfigurationFile, "Setting configuration file")
	assert.Equal(t, "/etc/rest2command/credentials.json", CredentialsFile, "Setting credentials")


}

func TestSetUpLog(t *testing.T) {

	levels := map[string]string {
		"info": "info",
		"debug": "debug",
		"panic": "panic",
		"error": "error",
		"warn": "warning",
		"": "info",
	}
	for key, value := range levels {
		os.Setenv("LOG_LEVEL", key)
		setUpLog()
		assert.Equal(t, log.GetLevel().String(), value, key + "OK")
	}

}