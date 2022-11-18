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

func (s *Server) GoTo(c *fiber.Ctx) error {
	alias := c.Query("alias")
	log.Infof("alias=%s", alias)
	var aliasDescriptors []*domain.AliasDescriptor
	if alias == "" {
		aliasDescriptors = s.aliasDescriptorRepository.List()
	} else {
		aliasDescriptors = s.aliasDescriptorRepository.ListByAlias(alias)
	}

	if 1 < len(aliasDescriptors) {
		return c.JSON(aliasDescriptors)
	}

	if len(aliasDescriptors) == 1 {
		destination := aliasDescriptors[0].Destination

		return c.Redirect(destination, fiber.StatusSeeOther)
	}

	// Not Found일 때는 google search
	return c.Redirect(fmt.Sprintf(defaultSearchURLFormat, alias), fiber.StatusSeeOther)
}
