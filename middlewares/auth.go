package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"ingenhouzs.com/chesshouzs/go-game/helpers"
	"ingenhouzs.com/chesshouzs/go-game/interfaces"
)

func Auth(r interfaces.Repository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				authHeader = c.QueryParam("sid")
				if authHeader == "" {
					return helpers.HttpResponse(c, http.StatusUnauthorized, "Unauthenticated", nil)
				}
			}
			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method: " + token.Header["alg"].(string))
				}
				return []byte(os.Getenv("JWT_SECRET")), nil
			})
			if err != nil {
				helpers.LogErrorCallStack(c, err)
				return helpers.HttpResponse(c, http.StatusUnauthorized, "Unauthenticated", nil)
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				userID := claims["id"].(string)
				user, err := r.GetUserDataByID(userID)
				if err != nil {
					return helpers.HttpResponse(c, http.StatusUnauthorized, "Unauthenticated", nil)
				}
				c.Set("user", user)
				return next(c)
			}

			return helpers.HttpResponse(c, http.StatusUnauthorized, "Unauthenticated", nil)
		}
	}
}
