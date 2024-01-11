package utils

import (
	"aytp/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitialiseDB() {
	var err error
	log.Print("Initialising Database...")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Could not find Database Url")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	log.Print("Successfully connected database!")

	setupModels(
		&models.Emailz{},
		&models.Admin{},
		&models.Patner{},
		&models.Student{},
		&models.Tutor{},
		&models.Course{},
		&models.Chapter{},
		&models.Section{},
		&models.Enroll{},
		&models.CourseAssign{},
		&models.Quiz{},
		&models.Question{},
		&models.QuizComplete{},
		&models.Meeting{},
		&models.File{},
		&models.Assignment{},
		&models.AssignmentFile{},
		&models.Submission{},
		&models.SubmissionFile{},
		&models.Badge{},
		&models.Event{},
		&models.Hero{},
		&models.PersonalMeetings{},
	)
}

func setupModels(models ...interface{}) {
	err := DB.AutoMigrate(models...)
	if err != nil {
		panic(err)
	}
}
