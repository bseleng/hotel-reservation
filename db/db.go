package db

const DBNAME = "hotel-reservation"
const DBURI = "mongodb://localhost:27017"
const TestDBNAME = "hotel-reservation-test"

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

type Pagination struct {
	Limit int64
	Page  int64
}
