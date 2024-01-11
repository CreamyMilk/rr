package models

import uuid "github.com/satori/go.uuid"



type Tutor struct {
	Base
	Name     string
	Email    string `gorm:"uniqueIndex"`
	Password string
	SchoolID     uuid.UUID
	RoomId       string
	SchoolName   string
	ProfilePhoto string
}
