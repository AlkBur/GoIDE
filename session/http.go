package session

import (
	"github.com/gorilla/sessions"
)

var HTTPSession *sessions.CookieStore

func init() {
	HTTPSession = sessions.NewCookieStore([]byte("GO_IDE"))
}
