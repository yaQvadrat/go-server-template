package httpapi

import (
	"app/internal/service"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

const (
	accIdEchoContext = "accId"
	roleEchoContext  = "role"
)

type AuthMiddleware struct {
	accountService service.Account
}

func (m *AuthMiddleware) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString, err := m.getBearerToken(c.Request())
		if err != nil {
			log.Error(fmt.Errorf("httpapi - AuthMiddleware.Auth: %w", err))
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		token, err := m.accountService.ParseToken(tokenString)
		if err != nil {
			log.Error(fmt.Errorf("httpapi - AuthMiddleware.Auth: %w", err))
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		c.Set(accIdEchoContext, token.AccId)
		c.Set(roleEchoContext, token.Role)
		return next(c)
	}
}

func (m *AuthMiddleware) getBearerToken(request *http.Request) (string, error) {
	const bearerPrefix = "Bearer "

	header := request.Header.Get(echo.HeaderAuthorization)

	if len(header) > len(bearerPrefix) && strings.EqualFold(bearerPrefix, header[:len(bearerPrefix)]) {
		return header[len(bearerPrefix):], nil
	}

	return "", ErrNotBearerToken
}
