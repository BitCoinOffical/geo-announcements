package models

import "time"

type Incident struct {
	Incident_id int
	Lat         float64
	Lon         float64
	Status      string
	Create_at   *time.Time
	Update_at   *time.Time
	Deleted_at  *time.Time
}
