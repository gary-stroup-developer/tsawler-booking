package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf will help prevent cross site attacts
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction, //in production this is true
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// this was just practice...saving for reference
// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("Hit the Page")
// 		next.ServeHTTP(w, r)
// 	})
// }
