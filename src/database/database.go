package database

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

type PdfFile struct {
	ID string `gorm:"primarykey;autoIncrement"` // index

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CompanyName string // 公司名
	CompanyLink string // 公司地址（项目动态） <TODO: 提取已受理时间，通过时间>
	ProcessData string // 审核过程时间：通过CompanyLink 爬虫提取
	Frequency   int    // 问询次数
	FileName    string // 下载的pdf名字

	QueryTextLen  int // 问询文本长度
	AnswerTextLen int // 回答文本长度

	TableLen      int // 回答文本中 图片表格占幅 <无法实现>
	AccountVocLen int // 回答文本中 会计词汇个数
	LawVocLen     int // 回答文本中 法律词汇个数
}

// 添加记录
func Add2db(pf PdfFile) {
	Db.Create(&pf) // 通过数据的指针来创建
}

func Updatedb(pf PdfFile) {
	Db.Model(&pf).Updates(pf)
}

func Init() {
	Db, _ = gorm.Open(sqlite.Open("./database/pdfFilesData.db"))
	Db.AutoMigrate(&PdfFile{})
}

func init() {
	Init()
}
