package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	dynamodbKeySeparator                    = ":"
	dynamodbAliasPartitionKeyPrefix         = "groups"
	dynamodbAliasSortKeyPrefix              = "aliases"
	dynamodbUserPartitionKeyPrefix          = "users"
	dynamodbAliasHitEventPartitionKeyPrefix = "hits"
)

type Alias struct {
	gorm.Model
	AliasGroup  string           `json:"aliasGroup" dynamodbav:"group"`
	Name        string           `json:"name"`
	Destination string           `json:"destination"`
	Hits        []*EventAliasHit `gorm:"foreignKey:alias_fk"`
}

func (alias Alias) GetDynamodbPartitionKey() string {
	return dynamodbAliasPartitionKeyPrefix + dynamodbKeySeparator + alias.AliasGroup
}

func (alias Alias) GetDynamodbSortKey() string {
	return dynamodbAliasSortKeyPrefix + dynamodbKeySeparator + alias.Name
}

func GetAliasGroupFrom(dynamodbPartitionKey string) string {
	prefix := dynamodbAliasPartitionKeyPrefix + dynamodbKeySeparator
	if !strings.HasPrefix(dynamodbPartitionKey, prefix) {
		return ""
	}

	return strings.TrimPrefix(dynamodbPartitionKey, prefix)
}

func GetAliasName(dynamodbSortKey string) string {
	prefix := dynamodbAliasSortKeyPrefix + dynamodbKeySeparator
	if !strings.HasPrefix(dynamodbSortKey, prefix) {
		return ""
	}

	return strings.TrimPrefix(dynamodbSortKey, prefix)
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

func (user User) GetDynamodbPartitionKey() string {
	return dynamodbUserPartitionKeyPrefix + dynamodbKeySeparator + user.Username
}

func (user User) GetDynamodbSortKey() string {
	return ""
}

func GetUserNameFrom(dynamodbPartitionKey string) string {
	prefix := dynamodbUserPartitionKeyPrefix + dynamodbKeySeparator
	if !strings.HasPrefix(dynamodbPartitionKey, prefix) {
		return ""
	}

	return strings.TrimPrefix(dynamodbPartitionKey, prefix)
}

type Group struct {
	Name      string  `gorm:"primaryKey"`
	Users     []*User `gorm:"many2many:users_groups;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
