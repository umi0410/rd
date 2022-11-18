package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) GoTo(c *fiber.Ctx) error {
	alias := c.Query("alias")
	if alias == "" {
		aliasDescriptors := s.aliasDescriptorRepository.List()
		return c.JSON(aliasDescriptors)
	}

	aliasDescriptors := s.aliasDescriptorRepository.ListByAlias(alias)
	if len(aliasDescriptors) != 1 {
		return c.JSON(aliasDescriptors)
	}

	destination := aliasDescriptors[0].Destination

	return c.Redirect(destination, fiber.StatusSeeOther)
}
