package debug

import (
	"github.com/semaphoreui/semaphore/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path"
	"runtime/pprof"
	"strconv"
	"time"
)

func Dump(w http.ResponseWriter, r *http.Request) {

	f, err := os.Create(path.Join(util.Config.Profiling.DumpsDir, "mem-"+strconv.Itoa(int(time.Now().Unix()))+".prof"))

	defer f.Close()

	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"context": "pprof",
		}).Error("error creating mem.prof")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = pprof.WriteHeapProfile(f)

	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"context": "pprof",
		}).Error("Failed to write memory profile")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
