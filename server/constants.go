package server

import (
	"github.com/go-while/go-cpu-mem-profiler"
	"github.com/go-while/nodare-db/logger"
)

var AVAIL_SUBDICKS = []uint32{10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000, 2000000000, 4000000000}

const DEFAULT_SUB_DICKS = 100

const DEFAULT_PW_LEN = 32 // admin/username:password
const DEFAULT_SUPERADMIN = "superadmin"

const DEFAULT_CONFIG_FILE = "config.toml"
const DATA_DIR = "dat"
const CONFIG_DIR = "cfg"

const DEFAULT_LOGS_FILE = "ndb.log"
const DEFAULT_LOGLEVEL_STR = "INFO"
const DEFAULT_LOGLEVEL_INT = ilog.INFO

const DEFAULT_SERVER_ADDR = "[::1]"
const DEFAULT_SERVER_TCP_PORT = "2420"
const DEFAULT_SERVER_UDP_PORT = "2240"
const DEFAULT_SERVER_SOCKET_PATH = "/tmp/ndb.socket"
const DEFAULT_SERVER_SOCKET_TCP_PORT = "3420"
const DEFAULT_SERVER_SOCKET_TLS_PORT = "4420"

const DEFAULT_TLS_PRIVKEY = "privkey.pem"
const DEFAULT_TLS_PUBCERT = "fullchain.pem"

// readline flags
const no_mode = 0x00
const modeADD = 0x11
const modeGET = 0x22
const modeSET = 0x33
const modeDEL = 0x44
const CaseAdded = 0x69
const CaseDupes = 0xB8
const CaseDeleted = 0x00
const CasePass = 0xFF
const FlagSearch = 0x42

// client proto flags
const magic1 = "1" // mem-prof
const magic2 = "2" // cpu-prof
const magicA = "A" // add
const magicD = "D" // del
const magicG = "G" // get
const magicL = "L" // list
const magicS = "S" // set
const magicZ = "Z" // quit

// socket proto flags
const KEY_LIMIT = 1024 * 1024 * 1024 // respond: CAN
const VAL_LIMIT = 1024 * 1024 * 1024 // respond: CAN
const EmptyStr = ""
const CR = "\r"
const LF = "\n"
const CRLF = CR + LF
const DOT = "."
const COM = ","
const SEM = ";"

// ASCII control characters
// [hex: 0 - 1F] // [DEC character code 0-31]
const NUL = string(0x00) // Null character 		// 0
const SOH = string(0x01) // Start of Heading 	// 1
const STX = string(0x02) // Start of Text 		// 2
const ETX = string(0x03) // End of Text 		// 3
const EOT = string(0x04) // End of Transmission // 4
const ENQ = string(0x05) // Enquiry 			// 5
const ACK = string(0x06) // Acknowledge 		// 6
const BEL = string(0x07) // Bell, Alert 		// 7
const SYN = string(0x16) // Synchronous Idle	// 22
const ETB = string(0x17) // End of Trans. Block // 23
const CAN = string(0x18) // Cancel 				// 24
const EOM = string(0x19) // End of medium 		// 25
const SUB = string(0x20) // Substitute  		// 26
const ESC = string(0x1B) // Escape 				// 27

// VIPER CONFIG DEFAULTS

const V_DEFAULT_SUB_DICKS = "100"
const V_DEFAULT_TLS_ENABLED = false
const V_DEFAULT_NET_WEBSRV_READ_TIMEOUT = 5
const V_DEFAULT_NET_WEBSRV_WRITE_TIMEOUT = 10
const V_DEFAULT_NET_WEBSRV_IDLE_TIMEOUT = 120

// VIPER CONFIG KEYS
const VK_ACCESS_SUPERADMIN_USER = "server.superadmin_user"
const VK_ACCESS_SUPERADMIN_PASS = "server.superadmin_pass"

const VK_LOG_LOGLEVEL = "log.loglevel"
const VK_LOG_LOGFILE = "log.logfile"

const VK_SETTINGS_BASE_DIR = "settings.base_dir"
const VK_SETTINGS_DATA_DIR = "settings.data_dir"
const VK_SETTINGS_SETTINGS_DIR = "settings.settings_dir"
const VK_SETTINGS_SUB_DICKS = "settings.sub_dicks"

const VK_SEC_TLS_ENABLED = "security.tls_enabled"
const VK_SEC_TLS_PRIVKEY = "security.tls_priv_key"
const VK_SEC_TLS_PUBCERT = "security.tls_pub_cert"

const VK_NET_WEBSRV_READ_TIMEOUT = "network.websrv_read_timeout"
const VK_NET_WEBSRV_WRITE_TIMEOUT = "network.websrv_write_timeout"
const VK_NET_WEBSRV_IDLE_TIMEOUT = "network.websrv_idle_timeout"

const VK_SERVER_HOST = "server.bindip"
const VK_SERVER_PORT_TCP = "server.port"
const VK_SERVER_PORT_UDP = "server.port_udp"
const VK_SERVER_SOCKET_PATH = "server.socket_path"
const VK_SERVER_SOCKET_PORT_TCP = "server.socket_tcpport"
const VK_SERVER_SOCKET_PORT_TLS = "server.socket_tlsport"

var Prof *prof.Profiler
