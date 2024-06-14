package main

import (
	"flag"
	"github.com/go-while/go-cpu-mem-profiler"
	"github.com/go-while/nodare-db/database"
	"github.com/go-while/nodare-db/logger"
	"github.com/go-while/nodare-db/server"
	"log"
)

const MODE = 1

var (
	flag_configfile   string
	flag_logfile  string
	flag_hashmode int

	Prof *prof.Profiler
	// setHASHER sets prefered hash algo
	// [ HASH_siphash | HASH_FNV32A | HASH_FNV64A ]
	// TODO config value HASHER
	setHASHER = database.HASH_FNV64A
)

func main() {
	Prof = prof.NewProf()
	server.Prof = Prof

	// capture the flags: overwrites config file settings!
	flag.StringVar(&flag_configfile, "config", server.DEFAULT_CONFIG_FILE, "path to config.toml")
	flag.IntVar(&flag_hashmode, "hashmode", database.HASH_FNV64A, "sets hashmode:\n sipHash = 1\n FNV32A = 2\n FNV64A = 3\n")
	flag.StringVar(&flag_logfile, "logfile", "", "path to ndb.log")
	flag.Parse()

	database.HASHER = flag_hashmode
	// this first line prints LOGLEVEL="XX" to console but will never showup in logfile!
	logs := ilog.NewLogger(ilog.GetEnvLOGLEVEL(), flag_logfile)
	cfg, sub_dicks := server.NewViperConf(flag_configfile, logs)

	switch MODE {
	case 0:

	case 1:
		db := database.NewDICK(logs, sub_dicks)
		if database.HASHER == database.HASH_siphash {
			db.XDICK.GenerateSALT()
		}
		srv := server.NewFactory().NewNDBServer(cfg, server.NewXNDBServer(db, logs), logs)
		logs.Debug("Mode 1: Loaded vcfg='%#v'", cfg)
		//suckDickCh <- sub_dicks // read sub_dicks from config, pass to suckDickCh so we can create subDICKs
		logs.Debug("Mode 1: Created DB sub_dicks=%d", sub_dicks)
		//<-waitCh
		logs.Debug("Mode 1: Booted sub_dicks=%d srv='%v'", sub_dicks, srv)
		//host := vcfg.GetString(VK_SERVER_HOST)
		//log.Printf("MAIN Debug host='%v'", host)
		//log.Printf("MAIN Debug srv='%#v'", srv)
		if logs.IfDebug() {
			logs.Debug("launching PprofWeb @ :1234")
			go Prof.PprofWeb(":1234")
		}
		srv.Start()
		srv.Stop()
	default:
		log.Fatalf("Invalid MODE=%d", MODE)
	}

} // end func main
