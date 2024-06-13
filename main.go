package main

import (
	"github.com/go-while/go-cpu-mem-profiler"
	"github.com/go-while/nodare-db/database"
	"github.com/go-while/nodare-db/logger"
	"github.com/go-while/nodare-db/server"
	"log"
)

const MODE = 1

var (
	Prof *prof.Profiler
	// setHASHER sets prefered hash algo
	// [ HASH_siphash | HASH_FNV32A | HASH_FNV64A ]
	// TODO config value HASHER
	setHASHER = database.HASH_FNV64A
)

func main() {
	Prof = prof.NewProf()
	server.Prof = Prof

	database.HASHER = setHASHER
	logs := ilog.NewLogger(ilog.GetEnvLOGLEVEL())
	sdCh := make(chan uint32, 1)  // buffered or deadlocks
	waitCh := make(chan struct{}) // unbuffered is fine here

	switch MODE {
	case 1:
		db := database.NewDICK(logs, sdCh, waitCh)
		if database.HASHER == database.HASH_siphash {
			db.XDICK.GenerateSALT()
		}
		ndbServer := server.NewXNDBServer(db, logs)
		srv, vcfg, sub_dicks := server.NewFactory().GetWebServer(ndbServer, logs)
		logs.Debug("Mode 1: Loaded vcfg='%#v'", vcfg)
		sdCh <- sub_dicks // read sub_dicks from config, pass to sdCh so we can create subDICKs
		logs.Debug("Mode 1: Created DB sub_dicks=%d", sub_dicks)
		<-waitCh
		logs.Debug("Mode 1: Booted sub_dicks=%d srv='%v'", sub_dicks, srv)
		//host := vcfg.GetString("server.host")
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
