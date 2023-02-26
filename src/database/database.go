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
	CompanyLink string // 公司地址（项目动态）      TODO: 提取已受理时间，通过时间
	ProcessData string // 审核过程时间：通过CompanyLink 爬虫提取
	Frequency   int    // 问询次数
	FileName    string // 下载的pdf名字

	QueCount int    // 问题个数
	QueText  string // 问题内容

	AllTextLen         int // 问询函文本总长度
	QueryTextLen       int // 问询文本长度
	AnswerTextLen      int // 回答文本长度
	AnswerClearTextLen int // 回答文本去掉停用词且分词之后的字数

	TableLen         int // 回答文本中 图片表格占幅 <无法实现>
	AccountVocNumber int // 回答文本中 会计词汇个数
	AccountVocLen    int // 回答文本中 会计词汇字数
	LawVocNumber     int // 回答文本中 法律词汇个数
	LawVocLen        int // 回答文本中 法律词汇字数
}

func initDB() {
	Db, _ = gorm.Open(sqlite.Open("./database/pdfFilesData.db"))
	Db.AutoMigrate(&PdfFile{})
}

func init() {
	initDB()
}
