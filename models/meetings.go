package models

import (
	"time"
)

type Meeting struct {
	Base
	Time        time.Time
	Title       string
	Description string
	RoomId      string
	SessionID   string
	Recording   string
	CourseId    string
	CourseName  string
}
