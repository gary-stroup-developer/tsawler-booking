package repository

import (
	"github.com/gary-stroup-developer/tsawler-booking/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(r models.RoomRestriction) error

	//SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error)
}
