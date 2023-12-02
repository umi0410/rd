package entity

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Alias struct {
	gorm.Model
	AliasGroup  string           `json:"aliasGroup" dynamodbav:"group"`
	Name        string           `json:"name"`
	Destination string           `json:"destination"`
	Hits        []*EventAliasHit `gorm:"foreignKey:alias_fk"`
}

func (alias Alias) Validate() error {
	if len(alias.AliasGroup) == 0 {
		return errors.New(fmt.Sprintf("required field \"%s\" is empty", "AliasGroup"))
	}
	if len(alias.Name) == 0 {
		return errors.New(fmt.Sprintf("required field \"%s\" is empty", "Name"))
	}
	if len(alias.Destination) == 0 {
		return errors.New(fmt.Sprintf("required field \"%s\" is empty", "Destination"))
	}
	return nil
}

type EventAliasHit struct {
	gorm.Model
	Hit     bool `json:"hit"`
	AliasFK uint `json:"aliasFk"`
	//Alias *Alias `json:"alias"`
	// TODO: To implement a user system.
	// The type of User is just string for right now.
	User string `json:"user"`
}

type User struct {
	Username  string   `gorm:"primaryKey"`
	Groups    []*Group `gorm:"many2many:users_groups;" dynamodbav:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Group struct {
	Name      string  `gorm:"primaryKey"`
	Users     []*User `gorm:"many2many:users_groups;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
