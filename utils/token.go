package utils

import (
	"errors"
	"net/http"

	jwt "github.com/golang-jwt/jwt/v4"
)

var secretKey string = "rshsrhjdtyrstrtsh"

type userClaims struct {
	Id string `json:"Id"`
	jwt.RegisteredClaims
}

func createToken(phone_number string, code string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims{
		Id: phone_number + "#" + code,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "psyweb",
		},
	})
	return token.SignedString([]byte(secretKey))
}

func verifyToken(tokenStr string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}
	if _, ok := token.Claims.(*userClaims); ok && token.Valid {
		return true, nil
	}
	return false, errors.New("invalid token")
}

func AuthenticateUserLogin(w http.ResponseWriter, phone_number string, code string) error {
	// 登录成功，创建Token
	token_str, err := createToken(phone_number, code)
	if err != nil {
		return err
	}
	// 设置手机号的cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "PhoneNumber",
		Value:    phone_number,
		HttpOnly: true,
	})
	// 设置Token的cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Value:    token_str,
		HttpOnly: true,
	})
	return nil
}

func AuthenticateStaffLogin(w http.ResponseWriter, id string, password string) error {
	// 登录成功，创建Token
	token_str, err := createToken(id, password)
	if err != nil {
		return err
	}
	// 设置Token的cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Token",
		Value:    token_str,
		HttpOnly: true,
	})
	return nil
}

func IsLogged(w http.ResponseWriter, r *http.Request) (bool, error) {
	cookie, err := r.Cookie("Token")
	if err != nil {
		return false, err
	}
	ret, err := verifyToken(cookie.Value)
	if err != nil {
		return false, err
	}
	return ret, nil
}
