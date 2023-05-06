package config

import (
	"os"
	"path/filepath"
	"testing"
)

func generateTestConfig(configContent string, testConfigFileName string) (string, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	testConfigFileNamePath := filepath.Join(workingDir, testConfigFileName)

	testConfigFile, err := os.Create(testConfigFileNamePath)

	if err != nil {
		return "", err
	}

	_, err = testConfigFile.WriteString(configContent)

	if err != nil {
		return "", err
	}

	err = testConfigFile.Close()

	if err != nil {
		return "", err
	}

	return testConfigFileNamePath, nil
}

func TestConfig_Load(t *testing.T) {
	configExample := `
		title = "Example Config"
		symbols = ["AAPL", "tsla", "GOOG ", ""]
	`

	testConfigFileName, err := generateTestConfig(configExample, "test_config_load.toml")

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = os.Remove(testConfigFileName)
	}()

	config := Load(testConfigFileName)

	if err != nil {
		t.Fatal(err)
	}

	if config.Title != "Example Config" {
		t.Errorf("Expected title to be \"Example Config\", got \"%s\"", config.Title)
	}

	if config.Symbols[0] != "AAPL" {
		t.Errorf("Expected first symbol to be AAPL, got \"%s\"", config.Symbols[0])
	}

	if config.Symbols[1] != "TSLA" {
		t.Errorf("Expected third symbol to be trimmed to TSLA, got \"%s\"", config.Symbols[1])
	}

	if config.Symbols[2] != "GOOG" {
		t.Errorf("Expected third symbol to be trimmed to GOOG, got \"%s\"", config.Symbols[2])
	}

	if len(config.Symbols) > 3 {
		t.Errorf("Expected empty symbol to be removed, got %d symbols", len(config.Symbols))
	}
}

func TestConfig_Load_validate_title(t *testing.T) {
	configExample := `
		title = ""
		symbols = ["AAPL", "tsla", "GOOG ", ""]
	`

	testConfigFileName, err := generateTestConfig(configExample, "test_config_load_validate_title.toml")

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = os.Remove(testConfigFileName)
	}()

	_ = Load(testConfigFileName)

	if err != nil && err.Error() != "title is required" {
		t.Errorf("Expected error to be \"title is required\", got \"%s\"", err.Error())
	}
}

func TestConfig_Load_validate_symbol(t *testing.T) {
	configExample := `
		title = "Example Config"
		symbols = ["", "", "", ""]
	`

	testConfigFileName, err := generateTestConfig(configExample, "test_config_load_validate_symbol.toml")

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = os.Remove(testConfigFileName)
	}()

	_ = Load(testConfigFileName)

	if err != nil && err.Error() != "symbols is required" {
		t.Errorf("Expected error to be \"symbols is required\", got \"%s\"", err.Error())
	}
}

func TestLoad_MissingConfigurationFile(t *testing.T) {
	defer func() {
		t.Log("The code panicked as expected")
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	_ = Load("does_not_exist.toml")
}
