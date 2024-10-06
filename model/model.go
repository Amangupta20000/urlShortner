package model

import "time"

// import (
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// )

type URL struct {
	// ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}
