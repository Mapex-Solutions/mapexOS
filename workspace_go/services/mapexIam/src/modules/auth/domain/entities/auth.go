package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type Auth struct {
	ID      model.ObjectId `bson:"_id"`
	Name    string         `bson:"name,omitempty"`
	Email   string         `bson:"email,omitempty"`
	Created time.Time      `bson:"created,omitempty"`
}

func NewAuth(name, email string) *Auth {
	return &Auth{Name: name, Email: email}
}
