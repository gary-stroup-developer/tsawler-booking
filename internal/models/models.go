package models

import "github.com/gary-stroup-developer/tsawler-booking/internal/forms"

// Template data holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Form      *forms.Form
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
}

type Reservation struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}
