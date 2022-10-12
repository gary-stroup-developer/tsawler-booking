package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gary-stroup-developer/tsawler-booking/pkg/config"
	"github.com/gary-stroup-developer/tsawler-booking/pkg/handlers"
	"github.com/gary-stroup-developer/tsawler-booking/pkg/render"
)

const portNumber = ":8080"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	// change to true when in production
	app.InProduction = false

	// set the session parameters
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction //in production should be true

	// store session in AppConfig
	app.Session = session

	// create the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	// store the data in AppConfig
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	log.Fatal(srv.ListenAndServe())
}
