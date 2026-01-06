package models

import "time"

type Incident struct {
	Incident_id int
	Lat         string
	Lon         string
	Create_at   *time.Time
	Update_at   *time.Time
}
