package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"rd/config"
	"rd/repository"
	"rd/service"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	promCli "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	app              *fiber.App
	redirectionCount = promauto.NewCounterVec(promCli.CounterOpts{
		Name: "rd_redirection_count",
		Help: "The total number of redirections.",
	}, []string{"group", "alias", "destination"})
)

func NewServer(aliasRepo repository.AliasRepository, aliasSvc service.AliasService, authSvc service.AuthService, host, port string, tlsConfig config.TlsConfig) (*Server, error) {
	s := &Server{
		aliasRepo: aliasRepo,
		aliasSvc:  aliasSvc,
		authSvc:   authSvc,
		host:      host,
		port:      port,
		tlsConfig: tlsConfig,
	}

	app = fiber.New()
	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use(favicon.New(favicon.Config{
		File: "./assets/favicon.ico",
	}))

	app.Get("/goto", s.GoTo)
	app.Post("/aliases", s.CreateAlias)
	app.Get("/aliases", s.ListAliases)

	return s, nil
}

type Server struct {
	aliasRepo repository.AliasRepository
	aliasSvc  service.AliasService
	authSvc   service.AuthService
	host      string
	port      string
	tlsConfig config.TlsConfig
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	log.Infof("Starting running on %s", addr)

	switch s.tlsConfig.Mode {
	case config.TlsModeMutual:
		if err := app.ListenMutualTLS(addr, s.tlsConfig.CertFile, s.tlsConfig.PrivateKeyFile, s.tlsConfig.ClientCaCertFile); err != nil {
			return errors.WithStack(err)
		}
	case config.TlsModeSimple:
		if err := app.ListenTLS(addr, s.tlsConfig.CertFile, s.tlsConfig.PrivateKeyFile); err != nil {
			return errors.WithStack(err)
		}
	default:
		log.Warningf("Unknown TLS mode: %s, so just run the server in the plain text mode", s.tlsConfig.Mode)
		fallthrough
	case config.TlsModePlaintext:
		if err := app.Listen(addr); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (s *Server) GetUser(ctx *fiber.Ctx) (string, bool) {
	auth := ctx.Query("auth")
	if len(auth) == 0 {
		return "", false
	}

	return auth, true
}
