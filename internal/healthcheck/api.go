package healthcheck

import (
	"os/exec"
	"errors"
	"github.com/gofiber/fiber/v2"
)

var appver string

// RegisterHandlers registers the handlers that perform healthchecks.
func RegisterHandlers(r *fiber.App, ver string) {
	appver = ver
	r.Get("/healthcheck", healthchecking)
}

// healthcheck responds to a healthcheck request.
func healthchecking(c *fiber.Ctx) error {
	// Exec: cat /proc/loadavg
	lsCmd := exec.Command("cat", "/proc/loadavg")
	lsOut, err := lsCmd.Output()
	if err != nil {
		return errors.New("Loadavg command dont respond")
	}
	return c.SendString("OK " + string(lsOut))
}
