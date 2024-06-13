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
const DEFAULT_LOGLEVEL_STR = "INFO"
const DEFAULT_LOGLEVEL_INT = ilog.INFO

const DEFAULT_SERVER_ADDR = "[::1]"
const DEFAULT_SERVER_TCP_PORT = "2420"
const DEFAULT_SERVER_UDP_PORT = "2240"
const DEFAULT_SERVER_SOCKET_PATH = "/tmp/nodare-db.socket"
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

var Prof *prof.Profiler
