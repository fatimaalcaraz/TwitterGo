package jwt

import (
	"errors"
	"strings"

	"github.com/fatimaalcaraz/TwitterGo/models"
	jwt "github.com/golang-jwt/jwt/v5"
)

var Email string
var IdUsuario string

func ProcesoToken(tk string, JWTSign string) (*models.Claim, bool, string, error) {
	miclave := []byte(JWTSign)
	var claims models.Claim

	splitToken := strings.Split(tk, "Bearer")
	if len(splitToken) != 2 {
		return &claims, false, string(""), errors.New("Formato de token  invalido")
	}
	//trimspace saca todos los espacios al inicio  y al final de una cadena de caracteres
	tk = strings.TrimSpace(splitToken[1])
	tkn, err := jwt.ParseWithClaims(tk, &claims, func(token *jwt.Token) (interface{}, error) {
		return miclave, nil
	})
	if err == nil {
		//reutina que chequea contra la BD
	}
	if !tkn.Valid {
		return &claims, false, string(""), errors.New("Token invalido")
	}
	return &claims, false, string(""), err
}
