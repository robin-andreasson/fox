package fox

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
)

type RefreshOptions struct {
	AccessToken     TokenOptions
	RefreshToken    TokenOptions
	RefreshFunction func(refreshobj any) (any, error) // retrieve payload you plan to store inside access token
	Cookie          CookieAttributes

	init bool
}

type TokenOptions struct {
	Secret string // string used to encode jwt tokens
	Exp    int    // Milliseconds until token expires
	Nbf    int    // Milliseconds until before activated
}

var refreshOpt RefreshOptions

const bearer = `^bearer\s*`

/*
creates a refresh middleware

will create refresh and access tokens

searches for the "authorization" header with bearer when it comes to access token
*/
func Refresh(options RefreshOptions) {
	if options.AccessToken.Secret == "" || options.RefreshToken.Secret == "" {
		log.Panic("zero value secret is not allowed")
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

	if len(bearer) > 1 {
		accesstoken := bearer[1]

		if payload, err := validateToken(accesstoken, refreshOpt.AccessToken.Secret); err == nil && payload != nil {
			c.Refresh = payload

			return
		}
	}

	if refreshCookie == "" {
		return
	}

	refreshpayload, err := validateToken(refreshCookie, refreshOpt.RefreshToken.Secret)

	if err != nil {
		return
	}

	newpayload, err := refreshOpt.RefreshFunction(refreshpayload)

	if err != nil {
		c.Error = append(c.Error, err)
		return
	}

	newaccesstoken, err := generateToken(newpayload, refreshOpt.AccessToken)

	if err != nil {
		c.Error = append(c.Error, err)
		return
	}

	c.Refresh = newpayload
	c.SetHeader("X-Fox-Access-Token", newaccesstoken)
}

func validateToken(tokenStr string, secret string) (any, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if len(claims) == 0 {
			return nil, errors.New("no claims")
		}

		return claims["data"], nil
	}

	return nil, nil
}

func generateToken(payload any, tokenopt TokenOptions) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": payload,
		"iat":  now.Unix(),
		"exp":  now.Add(time.Millisecond * time.Duration(tokenopt.Exp)).Unix(),
		"nbf":  now.Add(time.Millisecond * time.Duration(tokenopt.Nbf)).Unix(),
	})

	return token.SignedString([]byte(tokenopt.Secret))
}
