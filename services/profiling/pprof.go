package profiling

import (
	"fmt"
	"github.com/semaphoreui/semaphore/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/pprof"
)

func StartProfiling() {
	if !util.Config.Profiling.Enabled {
		return
	}

	mux := http.NewServeMux()

	// Register pprof handlers explicitly
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	go func() {
		log.WithFields(log.Fields{
			"context": "profiling",
			"port":    util.Config.Profiling.Port,
		}).Infof("pprof listening")
		err := http.ListenAndServe(fmt.Sprintf(":%d", util.Config.Profiling.Port), mux)
		log.WithError(err).WithFields(log.Fields{
			"context": "profiling",
		}).Info("pprof stopped")
		log.Println()
	}()

}
