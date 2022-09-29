package models

type AttractionWaitTime struct {
	ID           uint `gorm:"primaryKey"`
	Time         int64
	Attraction   Attraction
	AttractionID uint
	Status       string
	SingleRider  int16 `gorm:"default:null"`
	WaitTime     int16 `gorm:"default:null"`
}
