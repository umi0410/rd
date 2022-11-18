package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"rd/service"
)

var app *fiber.App

func NewServer(aliasDescriptorRepository service.AliasDescriptorRepository) (*Server, error) {
	s := &Server{
		aliasDescriptorRepository: aliasDescriptorRepository,
	}

	app = fiber.New()
	app.Get("/goto", s.GoTo)

	return s, nil
}

type Server struct {
	aliasDescriptorRepository service.AliasDescriptorRepository
}

func (s *Server) Run(host, port string) error {
	log.WithField("goto", fmt.Sprintf("http://%s:%s/goto", host, port)).Info("Start running")
	if err := app.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
