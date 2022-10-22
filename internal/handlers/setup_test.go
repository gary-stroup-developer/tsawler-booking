package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gary-stroup-developer/tsawler-booking/internal/config"
	"github.com/gary-stroup-developer/tsawler-booking/internal/models"
	"github.com/gary-stroup-developer/tsawler-booking/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "internal/templates"

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})

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
		log.Fatal("Cannot create template cache")
	}

	// store the data in AppConfig
	app.TemplateCache = tc
	app.UseCache = true

	repo := NewRepo(&app)
	NewHandlers(repo)
	render.NewTemplates(&app)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.General)
	mux.Get("/majors-suite", Repo.Major)
	mux.Get("/contact", Repo.Contact)

	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability-json", Repo.JsonAvailability)
	mux.Post("/search-availability", Repo.PostAvailability)

	mux.Get("/make-reservations", Repo.Reservation)
	mux.Post("/make-reservations", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

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

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ../../templates
	pages, _ := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", pathToTemplates))

	//range through all file ending with *.pages.tmpl
	for _, page := range pages {
		// name will be base of file passed in i.e. home.page.tmpl
		name := filepath.Base(page)
		//ts = template set. the current page *template.Template is binded to ts
		ts := template.Must(template.New(name).ParseFiles(page))

		// checks to see if there are layout templates
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		// the page may need the layout template so ParseGlob binds the files to the ts variable
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	//log.Println(myCache)
	return myCache, nil
}
