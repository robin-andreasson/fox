package fox

import (
	"regexp"

	"github.com/golang-jwt/jwt"
)

type RefreshOptions struct {
	Secret          string
	RefreshFunction func(refresh_obj any) any
	Cookie          CookieAttributes
}

const bearer = `^bearer:\s*`

func Refresh(options RefreshOptions) {

}

func handleRefresh(authorization string, refreshCookie string, refresh any) {
	if authorization == "" && refreshCookie == "" {
		return
	}

	bearer := regexp.MustCompile(bearer).Split(authorization, 1)

	if len(bearer) < 1 {
		return
	}

	accesstoken := bearer[1]

	jwt.Parse(accesstoken)

	//Verify the jwt token inside of the authorization header

	//if verification came out as valid then parse the token body and

}
