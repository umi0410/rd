package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
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

func NewServer(aliasRepo repository.AliasRepository, aliasSvc service.AliasService, authSvc service.AuthService) (*Server, error) {
	s := &Server{
		aliasRepo: aliasRepo,
		aliasSvc:  aliasSvc,
		authSvc:   authSvc,
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
}

func (s *Server) Run(host, port string) error {
	log.WithField("goto", fmt.Sprintf("http://%s:%s/goto", host, port)).Info("Start running")
	if err := app.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
		return errors.WithStack(err)
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
