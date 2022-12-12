package database

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID int `gorm:"primarykey;autoIncrement"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CompanyName   string
	FileName      string
	QueryTextLen  int
	AnswerTextLen int
	TableLen      int
	AccountVocLen int
	LawVocLen     int
}
