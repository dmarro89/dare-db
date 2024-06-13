package server

import (
	"github.com/go-while/go-cpu-mem-profiler"
	"github.com/go-while/nodare-db/logger"
)

const DEFAULT_PW_LEN = 32 // admin/username:password
const DEFAULT_ADMIN = "superadmin"

const DEFAULT_CONFIG_FILE = "config.toml"
const DATA_DIR = "dat"
const CONFIG_DIR = "cfg"

const DEFAULT_LOG_FILE = "db.log"
const DEFAULT_LOGLEVEL_STRING string = "INFO"
const DEFAULT_LOGLEVEL_INT = ilog.INFO

const DEFAULT_SERVER_ADDR = "[::1]"
const DEFAULT_SERVER_TCP_PORT = "2420"
const DEFAULT_SERVER_UDP_PORT = "2240"
const DEFAULT_SERVER_SOCKET_PATH = "/tmp/nodare-db.socket"
const DEFAULT_SERVER_SOCKET_TCP_PORT = "3420"
const DEFAULT_SERVER_SOCKET_TLS_PORT = "4420"

const DEFAULT_TLS_PRIVKEY = "privkey.pem"
const DEFAULT_TLS_PUBCERT = "fullchain.pem"

const EmptyStr = ""
const CaseAdded = 0x69
const CaseDupes = 0xB8
const CaseDeleted = 0x00
const CasePass = 0xFF
const FlagSearch = 0x42

const no_mode = 0x00
const modeADD = 0x11
const modeGET = 0x22
const modeSET = 0x33
const modeDEL = 0x44

var Prof *prof.Profiler
