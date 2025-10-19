package profiler

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func StartPprof() {
	go func() {
		log.Println("pprof listening on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatal(err)
		}
	}()
}
