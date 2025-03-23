package db

const (
	DBURI      = "mongodb://localhost:27017"
	DBNAME     = "hotel-reservation"
	TESTDBNAME = "hotel-reservation-test"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
