package server

import (
	"errors"
	"fmt"

	"os"
	"path/filepath"
	"strings"

	"github.com/go-while/nodare-db/utils"
	"github.com/spf13/viper"
	"log"
)

var AVAIL_SUBDICKS = []uint32{10, 100, 1000, 10000, 100000, 1000000}

type Config interface {
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetUint32(key string) uint32
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
	err := os.MkdirAll(dirPath, 0755)
	switch {
	case err == nil:
		log.Printf("Directory created successfully: %s", dirPath) // info
	case os.IsExist(err):
		log.Printf("Directory already exists: %s", dirPath) // debug
	default:
		log.Printf("Error creating directory: %v", err) // error
	}
}

func createDefaultConfigFile(c *viper.Viper, cfgFile string) {

	log.Printf("Creating default configuration file")
	adminuser := DEFAULT_ADMIN
	adminpass := utils.GenerateRandomString(DEFAULT_PW_LEN)

	c.SetDefault("server.host", DEFAULT_SERVER_ADDR_STR)
	c.SetDefault("server.port", DEFAULT_SERVER_TCP_PORT_STR)
	c.SetDefault("server.port_udp", DEFAULT_SERVER_UDP_PORT_STR)
	c.SetDefault("server.admin_user", adminuser)
	c.SetDefault("server.admin_password", adminpass)

	c.SetDefault("log.log_level", DEFAULT_LOGLEVEL_STRING)
	c.SetDefault("log.log_file", DEFAULT_LOG_FILE)

	c.SetDefault("settings.base_dir", ".")
	c.SetDefault("settings.data_dir", DATA_DIR)
	c.SetDefault("settings.settings_dir", CONFIG_DIR)
	c.SetDefault("settings.sub_dicks", "1000")

	c.SetDefault("security.tls_enabled", false)
	// /etc/letsencrypt/live/(sub.)domain.com/fullchain.pem
	c.SetDefault("security.tls_cert_private", filepath.Join(CONFIG_DIR, DEFAULT_TLS_PUBCERT))
	// /etc/letsencrypt/live/(sub.)domain.com/privkey.pem
	c.SetDefault("security.tls_cert_public", filepath.Join(CONFIG_DIR, DEFAULT_TLS_PRIVKEY))

	c.SetDefault("network.websrv_read_timeout", 5)
	c.SetDefault("network.websrv_write_timeout", 10)
	c.SetDefault("network.websrv_idle_timeout", 120)

	c.WriteConfigAs(cfgFile)

	fmt.Printf("\nIMPORTANT! Generated ADMIN credentials! \n admin login: '%s' password: %s\n\n", adminuser, adminpass)
}

func mappingEnvsToConfig() {

	mapsEnvsToConfig["server.host"] = "NDB_HOST"
	mapsEnvsToConfig["server.port"] = "NDB_PORT"
	mapsEnvsToConfig["server.admin_user"] = "NDB_USER"
	mapsEnvsToConfig["server.admin_password"] = "NDB_PASSWORD"

	mapsEnvsToConfig["log.log_level"] = "LOGLEVEL"
	mapsEnvsToConfig["log.log_file"] = "LOG_FILE"

	mapsEnvsToConfig["settings.base_dir"] = "NDB_BASE_DIR"
	mapsEnvsToConfig["settings.data_dir"] = "NDB_DATA_DIR"
	mapsEnvsToConfig["settings.settings_dir"] = "NDB_CONFIG_DIR"
	mapsEnvsToConfig["settings.sub_dicks"] = "NDB_SUB_DICKS"

	mapsEnvsToConfig["security.tls_enabled"] = "NDB_TLS_ENABLED"
	mapsEnvsToConfig["security.tls_cert_private"] = "NDB_TLS_KEY"
	mapsEnvsToConfig["security.tls_cert_public"] = "NDB_TLS_CRT"

	mapsEnvsToConfig["network.websrv_read_timeout"] = "NDB_WEBSRV_READ_TIMEOUT"
	mapsEnvsToConfig["network.websrv_write_timeout"] = "NDB_WEBSRV_WRITE_TIMEOUT"
	mapsEnvsToConfig["network.websrv_idle_timeout"] = "NDB_WEBSRV_IDLE_TIMEOUT"
}

func ReadConfigsFromEnvs(c *viper.Viper) {
	log.Printf("READ ENV VARS")
	for key, value := range mapsEnvsToConfig {
		valueFromEnv, ok := os.LookupEnv(value)
		if ok {
			log.Printf("GOT NEW ENV key='%s' value='%v'", key, valueFromEnv)
			c.Set(key, valueFromEnv)
		} else {
			log.Printf("NO ENV: key='%s' val='%s' !ok", key, valueFromEnv)
		}
	}
}

func initDB(c *viper.Viper) (sub_dicks uint32) {

	dbBaseDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error in getting current working directory: %v", err)
	}
	os.Setenv("NDB_BASE_DIR", dbBaseDir)

	createDirectory(filepath.Join(dbBaseDir, CONFIG_DIR))
	createDirectory(filepath.Join(dbBaseDir, c.GetString("settings.data_dir")))

	setSUBDICKS := c.GetUint32("settings.sub_dicks")
	for _, v := range AVAIL_SUBDICKS {
		if setSUBDICKS == v {
			sub_dicks = setSUBDICKS
			return
		}
	}
	// reached here we did not get a valid sub_dicks value from config
	// always return at least 10 so we don't fail
	log.Printf("WARN invalid sub_dicks value=%d !! defaulted to 1000", setSUBDICKS)
	return 1000
} // end func initDB

func PrintConfigsToConsole(c *viper.Viper) {
	fmt.Printf("Print all configs that were set\n")
	for key, _ := range mapsEnvsToConfig {
		fmt.Printf("Config value '%v': '%v'\n", key, c.Get(key))
	}
}

func NewConfiguration(cfgFile string) (Config, uint32) {
	v := viper.New()
	mappingEnvsToConfig()
	v.SetConfigType("toml")

	if len(strings.TrimSpace(cfgFile)) == 0 {
		log.Printf("No configuration file in '%s' was supplied. Using default value: %s", cfgFile, DEFAULT_CONFIG_FILE)
		cfgFile = DEFAULT_CONFIG_FILE
	}

	if !checkFileExists(cfgFile) {
		log.Printf("Configuration file does not exist: %s", cfgFile)
		createDefaultConfigFile(v, cfgFile)
	}

	log.Printf("Using configuration file: %s", cfgFile)

	v.SetConfigFile(cfgFile)

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic("Config file was not found")
		} else {
			panic("Config file was found, but another error was produced")
		}
	}

	ReadConfigsFromEnvs(v)
	sub_dicks := initDB(v)
	PrintConfigsToConsole(v)
	return &ViperConfig{v}, sub_dicks
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
