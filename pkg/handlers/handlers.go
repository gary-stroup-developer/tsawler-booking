package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gary-stroup-developer/tsawler-booking/pkg/config"
	"github.com/gary-stroup-developer/tsawler-booking/pkg/models"
	"github.com/gary-stroup-developer/tsawler-booking/pkg/render"
)

// Repo the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}

func (m *Repository) General(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Major(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.gohtml", &models.TemplateData{})
}

type jsonData struct {
	Ok      bool
	Message string
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{})
}

func (m *Repository) JsonAvailability(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(jsonData{
		Ok:      true,
		Message: "Sending data through this route to the console",
	})

	if err != nil {
		log.Println(err)
	}

	w.Write(data)
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start")
	end := r.FormValue("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}
