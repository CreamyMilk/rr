package models

import uuid "github.com/satori/go.uuid"

type Admin struct {
	Base
	Name         string
	SchoolName   string
	SchoolID     uuid.UUID
	SchoolRef    string
	ProfilePhoto string
	Email        string `gorm:"uniqueIndex"`
	Password     string
}

