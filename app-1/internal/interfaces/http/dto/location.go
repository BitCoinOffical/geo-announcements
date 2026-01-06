package dto

type LocationDTO struct {
	Lat string `json:"lat" binding:"required,lat"`
	Lon string `json:"lon" binding:"required,lon"`
}
