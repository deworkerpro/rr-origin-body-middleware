package origin

import (
	"bytes"
	"net/http"
	"io/ioutil"

	"github.com/roadrunner-server/http/v4/attributes"
	"go.uber.org/zap"
)

const name = "origin"

type Configurer interface {
	UnmarshalKey(name string, out any) error
	Has(name string) bool
}

type Logger interface {
	NamedLogger(name string) *zap.Logger
}

type Plugin struct {
	log *zap.Logger
	cfg Configurer
}

func (p *Plugin) Init(logger Logger, cfg Configurer) error {
	p.log = logger.NamedLogger(name)
	p.cfg = cfg
	return nil
}

func (p *Plugin) Weight() uint {
	return 1
}

func (p *Plugin) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.log.Fatal("Body read error")
			http.Error(w, "Body read error", http.StatusBadRequest)
			return
		}
		originBody := string(bodyBytes)
		r = attributes.Init(r)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		attributes.Set(r, "origin-body", originBody)
		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) Name() string {
	return name
}
