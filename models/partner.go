package models

import uuid "github.com/satori/go.uuid"

type Patner struct {
	Base
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	SchoolID     uuid.UUID
	SchoolName   string
	ProfilePhoto string
}
