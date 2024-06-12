package server

import "github.com/go-while/nodare-db/logger"

const DEFAULT_PW_LEN = 32 // admin/username:password
const DEFAULT_ADMIN = "superadmin"

const DEFAULT_CONFIG_FILE = "config.toml"
const DATA_DIR = "dat"
const CONFIG_DIR = "cfg"

const DEFAULT_LOG_FILE = "db.log"
const DEFAULT_LOGLEVEL_STRING string = "INFO"
const DEFAULT_LOGLEVEL_INT = ilog.INFO

const DEFAULT_SERVER_ADDR_STR = "[::1]"
const DEFAULT_SERVER_TCP_PORT_STR = "2420"
const DEFAULT_SERVER_UDP_PORT_STR = "2240"

const DEFAULT_TLS_PRIVKEY = "privkey.pem"
const DEFAULT_TLS_PUBCERT = "fullchain.pem"

