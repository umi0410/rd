package server

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"rd/domain"
	"rd/entity"
	"rd/service"
)

var (
	defaultSearchURLFormat = "https://google.com/search?q=%s"
)

var (
	labelAlias       = "alias"
	labelDestination = "destination"
)

func (s *Server) GoTo(c *fiber.Ctx) error {
	qGroup := c.Query("group")
	qAlias := c.Query("alias")
	user, _ := s.GetUser(c)

	destination := ""
	defer func() {
		labels := map[string]string{
			"group":       qGroup,
			"alias":       qAlias,
			"destination": destination,
		}
		//if destination != "" {
		//	labels["destination"] = destination
		//}
		redirectionCount.With(labels).Inc()
		log.Infof("Increased the counter.")
	}()
	log.Infof("qGroup=%s, qAlias=%s", qGroup, qAlias)

	if len(qGroup) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("group is a required argument")
	}

	aliases, err := s.aliasSvc.GoTo(c.Context(), user, qGroup, qAlias)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err)
	}

	if 1 < len(aliases) {
		return c.JSON(aliases)
	}

	if len(aliases) == 1 {
		destination = aliases[0].Destination

		return c.Redirect(destination, fiber.StatusSeeOther)
	}

	// Not Found일 때는 google search
	return c.Redirect(fmt.Sprintf(defaultSearchURLFormat, qAlias), fiber.StatusSeeOther)
}

func (s *Server) CreateAlias(c *fiber.Ctx) error {
	inputAlias := new(entity.Alias)
	if err := c.BodyParser(inputAlias); err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	alias, err := s.aliasRepo.Create(inputAlias)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString(err.Error())
	}

	return c.JSON(alias)
}

func (s *Server) ListAliases(c *fiber.Ctx) error {
	qGroup := c.Query("group")
	qAlias := c.Query("alias")
	qSort := service.Sort(c.Query("sort"))
	if err := qSort.Validate(); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err)
	}
	user, _ := s.GetUser(c)

	var (
		aliases []*domain.Alias
		err     error
	)

	if len(qGroup) != 0 && len(qAlias) != 0 {
		aliases, err = s.aliasSvc.ListByGroupAndAlias(c.Context(), user, qGroup, qAlias, qSort)
	} else if len(qGroup) != 0 {
		aliases, err = s.aliasSvc.ListByGroup(c.Context(), user, qGroup, qSort)
	} else {
		aliases, err = s.aliasSvc.List(c.Context(), user)
	}
	if err != nil {
		return err
	}

	return c.JSON(aliases)
}
