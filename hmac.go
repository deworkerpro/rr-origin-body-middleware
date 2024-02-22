package hmac

import (
	"net/http"

	"github.com/roadrunner-server/http/attributes"
	"go.uber.org/zap"
)

const name = "hmac"

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
		p.log.Info("i'm here")
		r = attributes.Init(r)
		attributes.Set(r, "origin", "body")
		next.ServeHTTP(w, r)
	})
}

func (p *Plugin) Name() string {
	return name
}
