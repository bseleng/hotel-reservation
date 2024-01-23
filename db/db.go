package db

import "os"

const MongoDbEnvName = "MONGO_DB_NAME"

var DBNAME = os.Getenv(MongoDbEnvName)

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
