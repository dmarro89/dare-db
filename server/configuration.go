package server

import (
	"errors"
	"fmt"

	"os"
	"path/filepath"
	"strings"

	"github.com/dmarro89/dare-db/logger"
	"github.com/dmarro89/dare-db/utils"

	"github.com/spf13/viper"
)

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	IsSet(key string) bool
}

type ViperConfig struct {
	*viper.Viper
}

var mapsEnvsToConfig map[string]string = make(map[string]string)

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func createDirectory(dirPath string) {
	_, err := os.Stat(dirPath)
	if err != nil && os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			logger.Error("Error creating directory:", err)
		} else {
			logger.Info("Directory created successfully:", dirPath)
		}
	} else if err != nil {
		logger.Error("Error:", err)
	} else {
		logger.Debug("Directory already exists: ", dirPath)
	}
}

func createDefaultConfigFile(c *viper.Viper, cfgFile string) {

	var passwordNew string = utils.GenerateRandomString(12)

	logger.Info("Creating default configuration file")

	c.SetDefault("server.host", "127.0.0.1")
	c.SetDefault("server.port", "2605")
	c.SetDefault("server.admin_user", "admin")
	c.SetDefault("server.admin_password", passwordNew)

	c.SetDefault("log.log_level", "INFO")
	c.SetDefault("log.log_file", "daredb.log")

	c.SetDefault("settings.data_dir", DATA_DIR)
	c.SetDefault("settings.settings_dir", SETTINGS_DIR)

	c.SetDefault("security.tls_enabled", false)
	c.SetDefault("security.cert_private", filepath.Join(SETTINGS_DIR, "cert_private"))
	c.SetDefault("security.cert_public", filepath.Join(SETTINGS_DIR, "cert_public"))

	c.WriteConfigAs(cfgFile)

	fmt.Printf("\nIMPORTANT! Generate default password for admin on initial start. Store it securely. Password: %v\n\n", passwordNew)
}

func mappingEnvsToConfig() {

	mapsEnvsToConfig["server.host"] = "DAREDB_HOST"
	mapsEnvsToConfig["server.port"] = "DAREDB_PORT"
	mapsEnvsToConfig["server.root_user_name"] = "DAREDB_USER"
	mapsEnvsToConfig["server.root_user_password"] = "DAREDB_PASSWORD"

	mapsEnvsToConfig["log.log_level"] = "DAREDB_LOG_LEVEL"
	mapsEnvsToConfig["log.log_file"] = "DAREDB_LOG_FILE"

	mapsEnvsToConfig["settings.data_dir"] = "DAREDB_DATA_DIR"
	mapsEnvsToConfig["settings.base_dir"] = "DAREDB_BASE_DIR"
	mapsEnvsToConfig["settings.settings_dir"] = "DAREDB_SETTINGS_DIR"

	mapsEnvsToConfig["security.tls_enabled"] = "DARE_TLS_ENABLED"
	mapsEnvsToConfig["security.cert_private"] = "DAREDB_CERT_PRIVATE"
	mapsEnvsToConfig["security.cert_public"] = "DAREDB_CERT_PUBLIC"
}

func reReadConfigsFromEnvs(c *viper.Viper) {
	logger.Info("Re-reading configurations from environmental variables")
	for key, value := range mapsEnvsToConfig {
		valueFromEnv, ok := os.LookupEnv(value)
		if ok {
			fmt.Println(key, valueFromEnv)
			c.Set(key, valueFromEnv)
		}
	}
}

func initDBDirectories(c *viper.Viper) {

	dbBaseDir, err := os.Getwd()
	if err != nil {
		logger.Error("Error in getting current working directory:", err)
	}
	os.Setenv("DAREDB_BASE_DIR", dbBaseDir)

	createDirectory(filepath.Join(dbBaseDir, SETTINGS_DIR))
	createDirectory(filepath.Join(dbBaseDir, c.GetString("settings.data_dir")))
}

func PrintConfigsToConsole() {
	fmt.Printf("Print all configs that were set\n")
	for key, _ := range mapsEnvsToConfig {
		fmt.Printf("Config value for for %v is: %v\n", key, viper.Get(key))
	}
}

func NewConfiguration(cfgFile string) Config {
	v := viper.New()
	mappingEnvsToConfig()
	v.SetConfigType("toml")

	if len(strings.TrimSpace(cfgFile)) == 0 {
		logger.Info("No configuration file was supplied. Using default value: ", DEFAULT_CONFIG_FILE)
		cfgFile = DEFAULT_CONFIG_FILE
	}

	isFileExist := checkFileExists(cfgFile)

	if !isFileExist {
		logger.Info("Configuration file does not exist: ", cfgFile)
		createDefaultConfigFile(v, cfgFile)
	}

	logger.Info("Using configuration file: ", cfgFile)

	v.SetConfigFile(cfgFile)

	if err := v.ReadInConfig(); err != nil {
		logger.Fatal("Error reading config file:", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file was not found")
		} else {
			panic("Config file was found, but another error was produced")
		}
	}

	reReadConfigsFromEnvs(v)
	initDBDirectories(v)
	return &ViperConfig{v}
}

func (c *ViperConfig) Get(key string) interface{} {
	return c.Viper.Get(key)
}

func (c *ViperConfig) GetString(key string) string {
	return c.Viper.GetString(key)
}

func (c *ViperConfig) GetBool(key string) bool {
	return c.Viper.GetBool(key)
}

func (c *ViperConfig) IsSet(key string) bool {
	return c.Viper.IsSet(key)
}
