package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"rd/domain"
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

	var aliases []*domain.Alias
	if len(qGroup) != 0 && len(qAlias) != 0 {
		aliases = s.aliasRepository.ListByGroupAndAlias(qGroup, qAlias)
	} else if qAlias == "" {
		aliases = s.aliasRepository.ListByGroup(qGroup)
	} else {
		aliases = s.aliasRepository.List()
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

func (s *Server) ListAliases(c *fiber.Ctx) error {
	return c.JSON(s.aliasRepository.List())
}
