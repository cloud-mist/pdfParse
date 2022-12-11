package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID int `gorm:"primarykey"`

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

func main() {
	// csvFilePath := "../material/company-file-data/company-file-all.csv"
	// download.ReadCsvAndDownLoad(csvFilePath)
}
