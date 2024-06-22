package server

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const TEST_CONFIG_FILE string = "config-test.toml"
const TEST_CONFIGURATION_DIR = "server"

var TEST_FOLDERS = []string{"data", "settings"}

func SetupTestConfiguration() Config {
	checkCorrectTestDirectory()
	return NewConfiguration(TEST_CONFIG_FILE)
}

func TeardownTestConfiguration() {
	checkCorrectTestDirectory()
	err := removeFileOrDirIfExists(TEST_CONFIG_FILE)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("File was removed successfully (if it existed):", TEST_CONFIG_FILE)
	}

	for _, folder := range TEST_FOLDERS {
		err := removeFileOrDirIfExists(folder)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("Folder was removed successfully (if it existed):", TEST_CONFIG_FILE)
		}
	}
}

// check, if tests run in the right directory
func checkCorrectTestDirectory() {
	baseDir, _ := os.Getwd()
	if !strings.HasSuffix(baseDir, TEST_CONFIGURATION_DIR) {
		panic("Wrong directory for running this test. Possibility to delete data and settings folders.")
	}
}

func removeFileOrDirIfExists(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = os.RemoveAll(filePath)
	if err != nil {
		return fmt.Errorf("failed to remove file/directory: %w", err)
	}
	return nil
}

func TestDefaultParameters(t *testing.T) {

	testConfig := SetupTestConfiguration()
	defer TeardownTestConfiguration()

	// Check if the values are correctly set
	assert.Equal(t, "127.0.0.1", testConfig.GetString("server.host"), "Host should be '127.0.0.1'")
	assert.Equal(t, "2605", testConfig.GetString("server.port"), "Port should be '2605'")
	assert.Equal(t, "admin", testConfig.GetString("server.admin_user"), "Admin name should be 'admin'")
	assert.Equal(t, "INFO", testConfig.GetString("log.log_level"), "Must be 'INFO'")
	assert.Equal(t, "daredb.log", testConfig.GetString("log.log_file"), "Must be 'daredb.log'")
	assert.Equal(t, false, testConfig.GetBool("security.tls_enabled"), "Must be 'false'")
}

func TestConfigurationConstants(t *testing.T) {

	testConfig := SetupTestConfiguration()
	defer TeardownTestConfiguration()

	// Check if the values are correctly set
	assert.Equal(t, "config.toml", DEFAULT_CONFIG_FILE, "Host should be 'config.toml'")
	assert.Equal(t, "data", DATA_DIR, "Host should be 'data'")
	assert.Equal(t, "settings", SETTINGS_DIR, "Host should be 'settings'")
	assert.Equal(t, SETTINGS_DIR, testConfig.GetString("settings.settings_dir"), "Host should be 'settings'")
}

func TestConfiguratioReReadFeature(t *testing.T) {
	t.Setenv("DARE_PORT", "2606")

	testConfig := SetupTestConfiguration()
	defer TeardownTestConfiguration()

	// Check if the values are correctly set
	assert.Equal(t, "2606", testConfig.GetString("server.port"), "Port should be '2606'")
}
