package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gary-stroup-developer/tsawler-booking/internal/config"
	"github.com/gary-stroup-developer/tsawler-booking/internal/driver"
	"github.com/gary-stroup-developer/tsawler-booking/internal/forms"
	"github.com/gary-stroup-developer/tsawler-booking/internal/helpers"
	"github.com/gary-stroup-developer/tsawler-booking/internal/models"
	"github.com/gary-stroup-developer/tsawler-booking/internal/render"
	"github.com/gary-stroup-developer/tsawler-booking/internal/repository"
	"github.com/gary-stroup-developer/tsawler-booking/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// func NewTestRepo(a *config.AppConfig) *Repository {
// 	return &Repository{
// 		App:a,
// 		DB: dbrepo.NewTestingRepo(a),
// 	}
// }

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.gohtml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "about.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.gohtml", &models.TemplateData{})
}

func (m *Repository) General(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.gohtml", &models.TemplateData{})
}

func (m *Repository) Major(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.gohtml", &models.TemplateData{})
}

type jsonData struct {
	Ok      bool
	Message string
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	data := make(map[string]interface{})
	data["reservation"] = models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))

	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	//will be passed back to make-reservation.page.gohtml template to populate the input values
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	// create a new form struct
	form := forms.New(r.PostForm)

	// check if all required fields have input. Check if first_name meets length requirement
	// checks if email is a valid email
	form.Required("first_name", "last_name", "email")
	form.MinLnegth("first_name", 3)
	form.IsEmail("email")

	// if form has errors, send back form info
	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form, //errors used to trigger message and bootstrap class
			Data: data, // reservation fields populate the input field values
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

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
	render.Template(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.FormValue("start")
	end := r.FormValue("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// remove session information since we obtain that in the reservation variable
	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data: data,
	})

}
