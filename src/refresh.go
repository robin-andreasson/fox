package fox

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
)

type RefreshOptions struct {
	Secret          string
	AccessToken     TokenOptions
	RefreshToken    TokenOptions
	RefreshFunction func(refreshobj any) (any, error) // retrieve payload you plan to store inside access token
	Cookie          CookieAttributes

	init bool
}

type TokenOptions struct {
	Exp int // Milliseconds until token expires
	Nbf int // Milliseconds until before activated
}

var refreshOpt RefreshOptions

const bearer = `^bearer:\s*`

func Refresh(options RefreshOptions) {
	if options.Secret == "" {
		log.Panic("zero value Secret is not allowed")
	}

	if options.RefreshFunction == nil {
		log.Panic("refresh function is required")
	}

	refreshOpt = options
	refreshOpt.init = true
}

func handleRefresh(authorization string, refreshCookie string, c *Context) {

	if !refreshOpt.init {
		return
	}

	if authorization == "" && refreshCookie == "" {
		return
	}

	bearer := regexp.MustCompile(bearer).Split(authorization, 2)

	if len(bearer) <= 1 {
		return
	}

	accesstoken := bearer[1]

	if payload, err := validateToken(accesstoken); err == nil {
		c.Refresh["payload"] = payload

		return
	} else {
		fmt.Println("ACCESS TOKEN ERROR: ", err)
	}

	if refreshCookie == "" {
		return
	}

	refreshpayload, err := validateToken(refreshCookie)

	if err != nil {
		return
	}

	newpayload, err := refreshOpt.RefreshFunction(refreshpayload)

	if err != nil {
		return
	}

	newaccesstoken, err := generateToken(newpayload, refreshOpt.AccessToken)

	if err != nil {
		return
	}

	c.Refresh["Payload"] = newpayload
	c.Refresh["Accesstoken"] = newaccesstoken
}

func validateToken(tokenStr string) (any, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(refreshOpt.Secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["data"], nil
	}

	return nil, err
}

func generateToken(payload any, tokenopt TokenOptions) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": payload,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Millisecond * time.Duration(tokenopt.Exp)).Unix(),
		"nbf":  now.Add(time.Millisecond * time.Duration(tokenopt.Nbf)).Unix(),
	})

	return token.SignedString([]byte(refreshOpt.Secret))
}
