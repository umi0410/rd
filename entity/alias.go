package entity

import "gorm.io/gorm"

type Alias struct {
	gorm.Model
	AliasGroup  string           `json:"alias_group"`
	Name        string           `json:"name"`
	Destination string           `json:"destination"`
	Hits        []*EventAliasHit `gorm:"foreignKey:alias_fk"`
}

type EventAliasHit struct {
	gorm.Model
	Hit     bool `json:"hit"`
	AliasFK int  `json:"alias_fk"`
	//Alias *Alias `json:"alias"`
	// TODO: To implement a user system.
	// The type of User is just string for right now.
	User string `json:"user"`
}
