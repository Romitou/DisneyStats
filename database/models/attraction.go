package models

type DisneyPark int

const (
	DisneylandPark DisneyPark = iota
	WaltDisneyStudios
	Unknown
)

type Attraction struct {
	ID       uint `gorm:"primaryKey"`
	EntityID string
	Name     string
	ParkID   DisneyPark
}
