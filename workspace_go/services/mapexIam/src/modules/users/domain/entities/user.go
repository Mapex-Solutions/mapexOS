package entities

import (
	"time"

	model "github.com/Mapex-Solutions/mapexGoKit/infrastructure/mongodb/model"
)

type AuthProvider struct {
	Type       string                 `bson:"type"`
	ExternalID *string                `bson:"externalId,omitempty"`
	Metadata   map[string]interface{} `bson:"metadata"`
}

type User struct {
	ID model.ObjectId `bson:"_id,omitempty"`

	Email                   string       `bson:"email"`
	Password                *string      `bson:"password,omitempty"`
	ChangePasswordNextLogin bool         `bson:"changePasswordNextLogin"`
	AuthProvider            AuthProvider `bson:"authProvider"`

	FirstName string  `bson:"firstName"`
	LastName  string  `bson:"lastName"`
	Phone     *string `bson:"phone"`
	JobTitle  *string `bson:"jobTitle,omitempty"`
	Enabled   bool    `bson:"enabled"`
	Avatar    *string `bson:"avatar,omitempty"`
	StartTour bool    `bson:"startTour"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

func (u *User) GetCreated() time.Time { return u.Created }
func (u *User) GetUpdated() time.Time { return u.Updated }
func (u *User) GetID() model.ObjectId { return u.ID }

// AuthProviderUpdateDTO allows patching sub-fields of AuthProvider.
// All fields are optional because this is for partial updates.
type AuthProviderUpdateDTO struct {
	Type       *string                 `bson:"type,omitempty"`
	ExternalID *string                 `bson:"externalId,omitempty"`
	Metadata   *map[string]interface{} `bson:"metadata,omitempty"`
}

// UserUpdateDTO is used for PATCH/UPDATE operations.
// Every field is optional (pointers), so nil means "ignore".
type UserUpdateDTO struct {
	Email                   *string                `bson:"email,omitempty"`
	Password                *string                `bson:"password,omitempty"`
	ChangePasswordNextLogin *bool                  `bson:"changePasswordNextLogin,omitempty"`
	AuthProvider            *AuthProviderUpdateDTO `bson:"authProvider,omitempty"`

	FirstName *string `bson:"firstName,omitempty"`
	LastName  *string `bson:"lastName,omitempty"`
	Phone     *string `bson:"phone,omitempty"`
	JobTitle  *string `bson:"jobTitle,omitempty"`
	Enabled   *bool   `bson:"enabled,omitempty"`
	Avatar    *string `bson:"avatar,omitempty"`
	StartTour *bool   `bson:"startTour,omitempty"`

	Created *time.Time `bson:"created"`
	Updated time.Time  `bson:"updated"`
}

func (u *UserUpdateDTO) GetCreated() *time.Time { return u.Created }
func (u *UserUpdateDTO) GetUpdated() time.Time  { return u.Updated }
