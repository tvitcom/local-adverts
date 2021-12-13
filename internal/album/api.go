package album

import (
	// "github.com/go-ozzo/ozzo-routing/v2"
	// "net/http"
	"github.com/gofiber/fiber/v2"
	"github.com/tvitcom/local-adverts/internal/errors"
	"github.com/tvitcom/local-adverts/pkg/log"
	"github.com/tvitcom/local-adverts/pkg/pagination"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, agregator Agregator, authHandler routing.Handler, logger log.Logger) {
	res := resource{agregator, logger}

	r.Get("/albums/<id>", res.get)
	r.Get("/albums", res.query)

	r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/albums", res.create)
	r.Put("/albums/<id>", res.update)
	r.Delete("/albums/<id>", res.delete)
}

type resource struct {
	agregator Agregator
	logger  log.Logger
}

func (res resource) get(c *routing.Context) error {
	album, err := res.agregator.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(album)
}

func (res resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	count, err := res.agregator.Count(ctx)
	if err != nil {
		return err
	}
	pages := pagination.NewFromRequest(c.Request, count)
	albums, err := res.agregator.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = albums
	return c.Write(pages)
}

func (res resource) create(c *routing.Context) error {
    //  p := new(Person)
    // if err := c.BodyParser(p); err != nil {
    //     return err
    // }
	var input CreateAlbumRequest
	if err := c.Read(&input); err != nil {
		res.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	album, err := res.agregator.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(album, http.StatusCreated)
}

func (res resource) update(c *routing.Context) error {
	var input UpdateAlbumRequest
	if err := c.Read(&input); err != nil {
		res.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	album, err := res.agregator.Update(c.Request.Context(), c.Param("id"), input)
	if err != nil {
		return err
	}

	return c.Write(album)
}

func (res resource) delete(c *routing.Context) error {
	album, err := res.agregator.Delete(c.Request.Context(), c.Param("id"))
	if err != nil {
		return err
	}

	return c.Write(album)
}

