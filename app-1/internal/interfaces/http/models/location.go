package models

import "time"

type Location struct {
	Location_id int
	User_id     string
	Lat         float64
	Lon         float64
	Status      string
	Create_at   *time.Time
}

type DangerousZones struct {
	Zone_id int
	Lat     float64
	Lon     float64
	Distant string
}

type UsersInDangerousZones struct {
	Zone_id   int
	UserCount int
}
