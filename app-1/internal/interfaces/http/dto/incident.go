package dto

type IncidentDTO struct {
	Lat float64 `json:"lat" binding:"required,lat"`
	Lon float64 `json:"lon" binding:"required,lon"`
}
