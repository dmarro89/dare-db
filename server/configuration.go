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
	viper            *viper.Viper
	logger           logger.Logger
	mapsEnvsToConfig map[string]string
}

func (c *ViperConfig) checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func (c *ViperConfig) createDirectory(dirPath string) {
	err := os.MkdirAll(dirPath, 0755)
	switch {
	case err == nil:
		c.logger.Info("Directory created successfully: ", dirPath)
	case os.IsExist(err):
		c.logger.Debug("Directory already exists: ", dirPath)
	default:
		c.logger.Error("Error creating directory:", err)
	}
}

func (c *ViperConfig) createDefaultConfigFile(cfgFile string) {
	var passwordNew string = utils.GenerateRandomString(12)

	c.logger.Info("Creating default configuration file")

	c.viper.SetDefault("server.host", "127.0.0.1")
	c.viper.SetDefault("server.port", "2605")
	c.viper.SetDefault("server.admin_user", "admin")
	c.viper.SetDefault("server.admin_password", passwordNew)

	c.viper.SetDefault("log.log_level", "INFO")
	c.viper.SetDefault("log.log_file", "daredb.log")

	c.viper.SetDefault("settings.data_dir", DATA_DIR)
	c.viper.SetDefault("settings.settings_dir", SETTINGS_DIR)

	c.viper.SetDefault("security.tls_enabled", false)
	c.viper.SetDefault("security.cert_private", filepath.Join(SETTINGS_DIR, "cert_private.pem"))
	c.viper.SetDefault("security.cert_public", filepath.Join(SETTINGS_DIR, "cert_public.pem"))

	c.viper.WriteConfigAs(cfgFile)

	fmt.Printf("\nIMPORTANT! Generate default password for admin on initial start. Store it securely. Password: %v\n\n", passwordNew)
}

func (c *ViperConfig) mappingEnvsToConfig() {
	c.mapsEnvsToConfig["server.host"] = "DARE_HOST"
	c.mapsEnvsToConfig["server.port"] = "DARE_PORT"
	c.mapsEnvsToConfig["server.admin_user"] = "DARE_USER"
	c.mapsEnvsToConfig["server.admin_password"] = "DARE_PASSWORD"

	c.mapsEnvsToConfig["log.log_level"] = "DARE_LOG_LEVEL"
	c.mapsEnvsToConfig["log.log_file"] = "DARE_LOG_FILE"

	c.mapsEnvsToConfig["settings.data_dir"] = "DARE_DATA_DIR"
	c.mapsEnvsToConfig["settings.base_dir"] = "DARE_BASE_DIR"
	c.mapsEnvsToConfig["settings.settings_dir"] = "DARE_SETTINGS_DIR"

	c.mapsEnvsToConfig["security.tls_enabled"] = "DARE_TLS_ENABLED"
	c.mapsEnvsToConfig["security.cert_private"] = "DARE_CERT_PRIVATE"
	c.mapsEnvsToConfig["security.cert_public"] = "DARE_CERT_PUBLIC"
}

func (c *ViperConfig) reReadConfigsFromEnvs() {
	c.logger.Info("Re-reading configurations from environmental variables")
	for key, value := range c.mapsEnvsToConfig {
		valueFromEnv, ok := os.LookupEnv(value)
		if ok {
			c.logger.Info("Use new configuration value from environmental variable for: ", key)
			fmt.Println(key, valueFromEnv)
			c.viper.Set(key, valueFromEnv)
		}
	}
}

func (c *ViperConfig) initDBDirectories() {
	dbBaseDir, err := os.Getwd()
	if err != nil {
		c.logger.Error("Error in getting current working directory:", err)
	}
	os.Setenv("DARE_BASE_DIR", dbBaseDir)

	c.createDirectory(filepath.Join(dbBaseDir, SETTINGS_DIR))
	c.createDirectory(filepath.Join(dbBaseDir, c.GetString("settings.data_dir")))
}

func (c *ViperConfig) PrintConfigsToConsole() {
	fmt.Printf("Print all configs that were set\n")
	for key := range c.mapsEnvsToConfig {
		fmt.Printf("Config value for for %v is: %v\n", key, c.Get(key))
	}
}

func NewConfiguration(cfgFile string) Config {
	logger := logger.NewDareLogger()
	if len(strings.TrimSpace(cfgFile)) == 0 {
		logger.Info("No configuration file was supplied. Using default value: ", DEFAULT_CONFIG_FILE)
		cfgFile = DEFAULT_CONFIG_FILE
	}

	v := viper.New()
	v.SetConfigType("toml")

	c := &ViperConfig{viper: viper.New(), logger: logger, mapsEnvsToConfig: make(map[string]string)}
	c.mappingEnvsToConfig()

	isFileExist := c.checkFileExists(cfgFile)

	if !isFileExist {
		c.logger.Info("Configuration file does not exist: ", cfgFile)
		c.createDefaultConfigFile(cfgFile)
	}

	logger.Info("Using configuration file: ", cfgFile)

	v.SetConfigFile(cfgFile)

	if err := v.ReadInConfig(); err != nil {
		c.logger.Fatal("Error reading config file:", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file was not found")
		} else {
			panic("Config file was found, but another error was produced")
		}
	}

	c.reReadConfigsFromEnvs()
	c.initDBDirectories()
	return &ViperConfig{viper: v, logger: logger, mapsEnvsToConfig: make(map[string]string)}
}

func (c *ViperConfig) Get(key string) interface{} {
	return c.viper.Get(key)
}

func (c *ViperConfig) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *ViperConfig) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *ViperConfig) IsSet(key string) bool {
	return c.viper.IsSet(key)
}
