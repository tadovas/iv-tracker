package rest

import (
	"flag"
	"net/http"
)

type Credentials struct {
	Username string
	Password string
}

func (c Credentials) IsSet() bool {
	return c.Username != ""
}

func RegisterFlags(c *Credentials) {
	flag.StringVar(&c.Username, "basic.username", "", "Username for basic auth user")
	flag.StringVar(&c.Password, "basic.password", "", "Password for basic auth user")
}

func BasicAuth(credentials Credentials) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, request *http.Request) {
			user, pass, ok := request.BasicAuth()
			if !ok {
				http.Error(writer, "Username or password missing", http.StatusUnauthorized)
				return
			}
			if user != credentials.Username && pass != credentials.Password {
				http.Error(writer, "Invalid username/password", http.StatusForbidden)
				return
			}
			next.ServeHTTP(writer, request)
		}
		return http.HandlerFunc(fn)
	}
}
