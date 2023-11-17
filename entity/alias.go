package entity

import (
	"time"

	"gorm.io/gorm"
)

type Alias struct {
	gorm.Model
	AliasGroup  string           `json:"aliasGroup"`
	Name        string           `json:"name"`
	Destination string           `json:"destination"`
	Hits        []*EventAliasHit `gorm:"foreignKey:alias_fk"`
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
	Groups    []*Group `gorm:"many2many:users_groups;"`
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
