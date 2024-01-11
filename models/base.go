package models

import (
	"aytp/graph/model"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var RoleModels = map[model.Role]interface{}{
	model.RoleAdmin:   Admin{},
	model.RolePatner:  Patner{},
	model.RoleStudent: Student{},
	model.RoleTutor:   Tutor{},
}

type Base struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
}

type UserInfo struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name         string
	Email        string `gorm:"uniqueIndex"`
	Password     string
	SchoolID     uuid.UUID
	SchoolName   string
	CreatedAt    time.Time `gorm:"not null"`
	ProfilePhoto string
}

func (u UserInfo) CreateToGraphData() *model.User {
	return &model.User{
		ID:           u.ID.String(),
		Name:         u.Name,
		Email:        u.Email,
		JoinedAt:     u.CreatedAt,
		ProfilePhoto: u.ProfilePhoto,
	}
}

func (base *Base) BeforeCreate(scope *gorm.DB) error {
	base.ID = uuid.NewV4()
	return nil
}
