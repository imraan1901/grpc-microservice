package db

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/imraan1901/grpc-microservice/internal/rocket"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

// New - returns a new store or error
func New() (Store, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbTable := os.Getenv("DB_TABLE")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSL_MODE")

	connectString := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		dbHost, dbPort, dbUsername, dbTable, dbPassword, dbSSLMode)

	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		return Store{}, err
	}

	return Store{
		db: db,
	}, nil

}

// GetRocketByID - retrieves a rocket from the database by id
func (s Store) GetRocketByID(id string) (rocket.Rocket, error) {
	var rkt rocket.Rocket
	row := s.db.QueryRow(
		`SELECT id, type, name FROM rockets where id=$1;`,
		id,
	)
	err := row.Scan(&rkt.ID, &rkt.Type, &rkt.Name)
	if err != nil {
		log.Print(err.Error())
		return rocket.Rocket{},nil
	}
	return rkt, nil
}

// InsertRocket - inserts a rocket into the rocket table
func (s Store) InsertRocket(rkt rocket.Rocket) (rocket.Rocket, error) {

	_, err := s.db.NamedQuery(
		`INSERT INTO rockets
		(id, name, type)
		VALUES (:id, :name, :type)`,
		rkt,
	)
	if err != nil {
		return rocket.Rocket{}, errors.New("failed to insert rocket into database")
	}
	return rocket.Rocket{
		ID: rkt.ID,
		Name: rkt.Name,
		Type: rkt.Type,
	}, nil
}

func (s Store) DeleteRocket(id string) error {
	return nil
}
