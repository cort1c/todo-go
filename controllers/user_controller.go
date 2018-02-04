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

type LoginRequest struct {
	Username string
	Password string
}

func LoginHandler(env *config.Env) echo.HandlerFunc {
	return echo.HandlerFunc(func(c echo.Context) error {
		req := &LoginRequest{}
		if err := c.Bind(req); err != nil {
			return err
		}
		user, err := models.FindUserByUsername(env.DB, req.Username)
		if err != nil {
			return err
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return echo.ErrUnauthorized
		}
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["name"] = user.Username
		claims["admin"] = false
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		t, err := token.SignedString([]byte(config.Secret))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	})
}
