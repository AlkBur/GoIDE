package session

import (
	"github.com/gorilla/sessions"
	"math/rand"
	"strconv"
	"time"
)

var HTTPSession *sessions.CookieStore

func init() {
	HTTPSession = sessions.NewCookieStore([]byte("GO_IDE"))
}

func GenId() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Int())
}
