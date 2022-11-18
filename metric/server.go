package metric

import (
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func Run() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":18090", nil); err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
}
