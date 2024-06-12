package main

import (
	"fmt"
	"github.com/go-while/nodare-db/client/clilib"
	"log"
	"sync"
	"time"
	"os"
)

var (
	wg sync.WaitGroup
)

func main() {
	stopChan := make(chan struct{}, 1)
	testWorker := true // runs a test after connecting
	daemon := false
	ssl := false
	mode := 1
	addr := "localhost:2420"

	netcli, err := client.NewClient(&client.Options{
		SSL:        ssl,
		Addr:       addr,
		Mode:       mode,
		Stop:       stopChan,
		Daemon:     daemon,
		TestWorker: testWorker,
	})
	time.Sleep(time.Second)
	if netcli == nil || err != nil {
		log.Printf("ERROR netcli='%v' err='%v'", netcli, err)
		return
	}
	//log.Printf("netcli='%#v'", netcli)
	parallel := 4

	items := 100000
	rounds := 10
	//testmap := make(map[string]string)
	parchan := make(chan struct{}, parallel)
	retchan := make(chan map[string]string, 1)
	log.Printf("starting insert")
	start := time.Now().Unix()

	// launch insert worker
	for r := 1; r <= rounds; r++ {
		time.Sleep(100*time.Millisecond)
		go func(r int, items int, parchan chan struct{}, retchan chan map[string]string) {
			parchan <- struct{}{} // locks
			testmap := make(map[string]string)
			log.Printf("launch insert worker r=%d", r)
			for i := 1; i <= items; i++ {
				// %010 leftpads i and r with 10 zeroes, like 17 => 0000000017
				key := fmt.Sprintf("atestKey%010d-r-%010d", i, r)
				val := fmt.Sprintf("atestVal%010d-r-%010d", i, r)
				res, err := netcli.Set(key, val)
				if err != nil {
					log.Fatalf("ERROR set key='%s' => val='%s' err='%v' res='%v'", key, val, err, res)
				}
				testmap[key] = val
			}
			<- parchan // returns lock
			retchan <- testmap
			log.Printf("returned insert worker r=%d set=%d", r, len(testmap))
		}(r, items, parchan, retchan)
	}

	log.Printf("wait for insert worker to return maps to test K:V")
	var capturemaps []map[string]string
	forever:
	for {
		select {
			case testmap := <- retchan:
				capturemaps = append(capturemaps, testmap)
				log.Printf("wait got a testmap have=%d want=%d", len(capturemaps), rounds)
			default:
				if len(capturemaps) == rounds {
					log.Printf("OK all testmaps returned, checking now...")
					break forever
				}
				time.Sleep(time.Millisecond*100)
		}
	}
	insert_end := time.Now().Unix()
	log.Printf("insert finished. checking... took %d sec", insert_end - start)

	// check all testmaps
	retint := make(chan int, len(capturemaps))
	for _, testmap := range capturemaps {
		go func(parchan chan struct{}, retint chan int, testmap map[string]string){
			parchan <- struct{}{} // locks
			checked := 0
			for k, v := range testmap {
				val, err := netcli.Get(k) // check GET
				if err != nil {
					log.Fatalf("ERROR get k='%s' err='%v'", k, err)
				}
				if val != v {
					log.Fatalf("ERROR verify k='%s' v='%s' != val='%s'", k, v, val)
					os.Exit(1)
				}
				checked++
			}
			<- parchan // returns lock
			retint <- checked
		}(parchan, retint, testmap)
	}


	checked := 0
	final:
	for {
		select {
			case aint := <- retint:
				checked += aint
				if checked == items*rounds {
					break final
				}
		}
	}
	test_end := time.Now().Unix()

	log.Printf("\n test parallel=%d total=%d/%d \n items/round=%d rounds=%d\n insert took %d sec \n check took %d sec \n total %d sec", parallel, checked, items*rounds, items, rounds, insert_end-start, test_end - insert_end, test_end - start)

	log.Printf("infinite wait on stopChan")
	<- stopChan
}
