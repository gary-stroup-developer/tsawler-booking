package dbrepo

import (
	"database/sql"

	"github.com/gary-stroup-developer/tsawler-booking/internal/config"
	"github.com/gary-stroup-developer/tsawler-booking/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}
