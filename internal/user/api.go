package user

import (
	"github.com/go-ozzo/ozzo-routing/v2"
	"github.com/tvitcom/local-adverts/internal/errors"
	"github.com/tvitcom/local-adverts/pkg/log"
	"github.com/tvitcom/local-adverts/pkg/pagination"
	"net/http"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, agregator Agregator, authHandler routing.Handler, logger log.Logger) {
	res := resource{agregator, logger}

	r.Get("/users/<id>", res.get)
	r.Get("/users", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/users", res.create)
	r.Put("/users/<id>", res.update)
	r.Delete("/users/<id>", res.delete)
}

type resource struct {
	agregator Agregator
	logger  log.Logger
}

func (res resource) get(c *routing.Context) error {
	user, err := res.agregator.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(user)
}

func (res resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := res.agregator.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	users, err := res.agregator.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = users
	return c.Write(pages)
}

func (res resource) create(c *routing.Context) error {
	var input CreateUserRequest
	if err := c.Read(&input); err != nil {
		res.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	user, err := res.agregator.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(user, http.StatusCreated)
}

func (res resource) update(c *routing.Context) error {
	var input UpdateUserRequest
	if err := c.Read(&input); err != nil {
		res.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	user, err := res.agregator.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(user)
}

func (res resource) delete(c *routing.Context) error {
	user, err := res.agregator.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(user)
}

