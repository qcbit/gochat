package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("auth")
	if err == http.ErrNoCookie {
		// not authenticated
		resp.Header().Set("Location", "/login")
		resp.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// some other error
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	// success - call the next handler
	h.next.ServeHTTP(resp, req)
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the 3rd-party login process
// format: /auth/{action}/{provider}
func loginHandler(resp http.ResponseWriter, req *http.Request) {
	segs := strings.Split(req.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(resp, fmt.Sprintf("Error when trying to get provider %s: %s",
				provider, err), http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(resp, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s:%s",
				provider, err), http.StatusInternalServerError)
			return
		}
		resp.Header.Set("Location", loginUrl)
		resp.WriteHeader(http.StatusTemporaryRedirect)
	default:
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(resp, "Auth action %s not support", action)
	}
}
