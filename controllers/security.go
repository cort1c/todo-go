package controllers

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dirges/todo/config"
	"github.com/dirges/todo/models"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

const secret = "mdcXkEe8paw4pANibwwcfuocKzM4gqBAKwExBhVkanEPPLtnhu2QTQJ7TVtkNXaV"

func LoginHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		user, err := models.GetUserByUsername(env.DB, username)
		if err != nil {
			return err
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return echo.ErrUnauthorized
		}
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = user.Username
		claims["admin"] = false
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte(secret))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	})
}
