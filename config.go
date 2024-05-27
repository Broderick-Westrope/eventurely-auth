package main

import (
	"errors"
	"os"
	"strings"
)

type configuration struct {
	addr        string
	supertokens *supertokensConfig
}

type supertokensConfig struct {
	adminEmails   []string
	connectionURI string
	apiKey        string
	appName       string
	apiDomain     string
	websiteDomain string
}

func loadConfig() (*configuration, error) {
	addr, err := getRequiredEnv("ADDR")
	if err != nil {
		return nil, err
	}
	config := &configuration{
		addr: addr,
	}

	prefix := "SUPERTOKENS__"
	{
		adminEmails, err := getSliceFromEnv(prefix+"ADMIN_EMAILS", ",")
		if err != nil {
			return nil, err
		}

		connectionURI, err := getRequiredEnv(prefix + "CONNECTION_URI")
		if err != nil {
			return nil, err
		}

		apikey, err := getRequiredEnv(prefix + "API_KEY")
		if err != nil {
			return nil, err
		}

		appName, err := getRequiredEnv(prefix + "APP_NAME")
		if err != nil {
			return nil, err
		}

		apiDomain, err := getRequiredEnv(prefix + "API_DOMAIN")
		if err != nil {
			return nil, err
		}

		websiteDomain, err := getRequiredEnv(prefix + "WEBSITE_DOMAIN")
		if err != nil {
			return nil, err
		}

		config.supertokens = &supertokensConfig{
			adminEmails:   adminEmails,
			connectionURI: connectionURI,
			apiKey:        apikey,
			appName:       appName,
			apiDomain:     apiDomain,
			websiteDomain: websiteDomain,
		}
	}

	return config, nil
}

func getRequiredEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", errors.New("missing value " + key)
	}
	return value, nil
}

func getSliceFromEnv(key, delim string) ([]string, error) {
	value, err := getRequiredEnv(key)
	if err != nil {
		return nil, err
	}

	return strings.Split(value, delim), nil
}
