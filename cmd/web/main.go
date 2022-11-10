package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gary-stroup-developer/tsawler-booking/internal/config"
	"github.com/gary-stroup-developer/tsawler-booking/internal/driver"
	"github.com/gary-stroup-developer/tsawler-booking/internal/handlers"
	"github.com/gary-stroup-developer/tsawler-booking/internal/models"
	"github.com/gary-stroup-developer/tsawler-booking/internal/render"
	"github.com/joho/godotenv"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	fmt.Printf("Starting application on port %s", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	log.Fatal(srv.ListenAndServe())
}

func run() (*driver.DB, error) {
	// Load the .env file in the current directory
	godotenv.Load()

	// or

	godotenv.Load(".env")

	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

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

	// connect to database
	pass := os.Getenv("DBPASSWORD")
	db, err := driver.ConnectSQL(fmt.Sprintf("host=localhost port=5432 dbname=bookings user=postgres password=%s", pass))

	if err != nil {
		log.Fatal("cannot connect to database! Dying...")
	}

	// create the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	// store the data in AppConfig
	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)

	return db, nil
}
