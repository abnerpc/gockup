package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Config records of current configuration
type Config struct {
	Target  string   `json:"target"`
	Sources []string `json:"sources"`
}

// Backuper records information to use on backup process
type Backuper struct {
	config *Config
}

// getAppConfigPath returns default app path in User home
func getAppConfigPath() string {
	homeDir, error := os.UserHomeDir()
	if error != nil {
		homeDir = "."
	}

	homeDir += "/.config/gockup"

	return homeDir
}

// readConfig try to read the configuration from default path
func readConfig() (*Config, error) {
	appConfigPath := getAppConfigPath()
	configFilePath := appConfigPath + "/config.json"

	configFileData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		os.MkdirAll(appConfigPath, os.ModePerm)
		err := ioutil.WriteFile(configFilePath, []byte("{}"), 0644)
		if err != nil {
			return nil, fmt.Errorf("Error creating config file at path: %s", configFilePath)
		}

		configFileData, err = ioutil.ReadFile(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("Error creating empty config file")
		}
	}

	var config Config
	err = json.Unmarshal(configFileData, &config)
	if err != nil {
		return nil, fmt.Errorf("Error parsing json file")
	}

	return &config, nil
}

func (c *Config) isValid() error {
	if len(c.Target) == 0 {
		return errors.New("Target not found")
	}

	if len(c.Sources) == 0 {
		return errors.New("Sources not found")
	}

	return nil
}

func main() {
	setTargetCmd := flag.NewFlagSet("set-target", flag.ExitOnError)
	targetPath := setTargetCmd.String("path", ".", "Path")
	addSourceCmd := flag.NewFlagSet("add-source", flag.ExitOnError)
	sourcePath := addSourceCmd.String("path", ".", "Path")

	switch os.Args[1] {

	case "set-target":
		setTargetCmd.Parse(os.Args[2:])
		path, _ := filepath.Abs(*targetPath)
		fmt.Print(path)
	case "add-source":
		addSourceCmd.Parse(os.Args[2:])
		path, _ := filepath.Abs(*sourcePath)
		fmt.Print(path)
	case "run":
		config, err := readConfig()
		if err != nil {
			log.Fatal(err)
		}

		if err := config.isValid(); err != nil {
			fmt.Printf("Invalid config: %s\n", err)
			return
		}

		fmt.Print("Valid\n")
	}

}
