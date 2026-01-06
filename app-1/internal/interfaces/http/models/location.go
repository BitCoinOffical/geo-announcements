package models

import "time"

type Location struct {
	Location_id int
	User_id     string
	Lat         string
	Lon         string
	Status      string
	Create_at   *time.Time
}
