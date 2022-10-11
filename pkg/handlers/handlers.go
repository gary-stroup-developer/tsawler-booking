package handlers

import (
	"net/http"
	"strconv"

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
	render.RenderTemplate(w, "home.page.gohtml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := map[string]string{
		"a": "Hello, again",
		"b": "Gary",
		"t": "Stroup",
		"z": strconv.FormatBool(Repo.App.UseCache),
	}
	render.RenderTemplate(w, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}
