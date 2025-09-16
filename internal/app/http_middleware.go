package app

import (
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/tamboto2000/otaqku-tasks/internal/common"
	"github.com/tamboto2000/otaqku-tasks/internal/modules/auth"
)

func AuthMiddleware(authSvc auth.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenStr := getTokenFromBearer(c)
			accId, err := authSvc.ValidateAccessToken(c.Request().Context(), tokenStr)
			if err != nil {
				return common.ErrorResponse(c, err)
			}

			c.Set("account_id", accId)

			return next(c)
		}
	}
}

var (
	bearerTokenRegex     = regexp.MustCompile(`^Bearer (?P<token>\S+)$`)
	tokenRegexGroupIndex = bearerTokenRegex.SubexpIndex("token")
)

func getTokenFromBearer(ectx echo.Context) string {
	matches := bearerTokenRegex.FindStringSubmatch(ectx.Request().Header.Get("Authorization"))
	if len(matches) == 0 {
		return ""
	}

	token := matches[tokenRegexGroupIndex]

	return token
}
