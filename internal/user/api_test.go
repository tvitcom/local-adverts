package user

import (
	"github.com/tvitcom/local-adverts/internal/auth"
	"github.com/tvitcom/local-adverts/internal/entity"
	"github.com/tvitcom/local-adverts/internal/test"
	"github.com/tvitcom/local-adverts/pkg/log"
	"net/http"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.User{
		{"123", "users123", time.Now(), time.Now()},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/userss", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/userss/123", "", nil, http.StatusOK, `*users123*`},
		{"get unknown", "GET", "/userss/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/userss", `{"name":"test"}`, header, http.StatusCreated, "*test*"},
		{"create ok count", "GET", "/userss", "", nil, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/userss", `{"name":"test"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/userss", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/userss/123", `{"name":"usersxyz"}`, header, http.StatusOK, "*usersxyz*"},
		{"update verify", "GET", "/userss/123", "", nil, http.StatusOK, `*usersxyz*`},
		{"update auth error", "PUT", "/userss/123", `{"name":"usersxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/userss/123", `"name":"usersxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/userss/123", ``, header, http.StatusOK, "*usersxyz*"},
		{"delete verify", "DELETE", "/userss/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/userss/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
