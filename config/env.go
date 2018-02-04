package config

import (
	"database/sql"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

const Secret string = "mdcXkEe8paw4pANibwwcfuocKzM4gqBAKwExBhVkanEPPLtnhu2QTQJ7TVtkNXaV"

type Env struct {
	DB *sql.DB
}

func (env *Env) GetCurrentUserID(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return int(claims["id"].(float64))
}
