package main

import (
	"fmt"
	"hello/download"
	"hello/parsepdf"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func main() {
	// mainDownload()
	// mainCalc()

}

func pdfToTxt() {
	pdfFiles := getPdfFilePath()
	for _, pdfFilePath := range pdfFiles {
		reg := regexp.MustCompile(`\d+-.*`)
		res := reg.FindString(pdfFilePath)
		fmt.Println(res)
		textFilePath := "../../txts/" + strings.ReplaceAll(res, "pdf", "txt")
		fmt.Println(textFilePath)
		parsepdf.ReadPdfV2(pdfFilePath, textFilePath)
	}
}

//1. 下载
func mainDownload() {
	csvFilePath := "../material/company-file-data/company-file-V3.csv"
	download.ReadCsvAndDownLoad(csvFilePath)
	csvFilePath = "../material/company-file-data/company-file-add.csv"
	download.ReadCsvAndDownLoad(csvFilePath)
}

//2. 统计
func mainCalc() {
	// pdfFiles := getPdfFilePath()
	// for _, pdfFilePath := range pdfFiles {
	pdfFilePath := "../../downloadsPDF/1-688165.SH_埃夫特-U_四轮反馈回复.pdf"
	fmt.Println("---------------------------------------------------------------------------")
	id := getPdfId(pdfFilePath)
	// 解析
	PrintPartDivMeg("解析部分开始")
	fmt.Printf("FileName:[%s]\n", pdfFilePath)
	fmt.Printf("<%v 解析: 分割成两部分>\n", getTime())
	// parsepdf.ReadPdfV2(pdfFilePath)
	parsepdf.DivideTwoParts()

	// 解析结果写入数据库
	fmt.Printf("<%v 解析：解析部分结果写入数据库>\n", getTime())
	parsepdf.WriteSomeParseResToDB(id)

	// 分词
	PrintPartDivMeg("分词部分开始")
	parsepdf.Divide("../material/tmpFile/tmpPart2.txt")

	PrintPartDivMeg("法律词汇统计")
	lawWordsFilePath := "../material/wordsFiles/law-words.txt"
	parsepdf.AddCompareWords(lawWordsFilePath)
	parsepdf.Count()
	fmt.Printf("<%v 分词：法律词汇结果写入数据库>\n", getTime())
	parsepdf.WriteWordsVocNum(id, "law")

	PrintPartDivMeg("会计词汇统计")
	accountWordsFilePath := "../material/wordsFiles/accounting-words.txt"
	parsepdf.AddCompareWords(accountWordsFilePath)
	parsepdf.Count()
	fmt.Printf("<%v 分词：会计词汇结果写入数据库>\n", getTime())
	parsepdf.WriteWordsVocNum(id, "account")
	fmt.Println("全部任务结束！")
}

// ------------------ 辅助函数 ----------------------------------------
func getTime() string {
	now := time.Now()
	hour := now.Hour()
	min := now.Minute()
	sec := now.Second()
	return fmt.Sprintf("%02d:%02d:%02d", hour, min, sec)
}

func PrintPartDivMeg(str string) {
	fmt.Println("********************************************")
	fmt.Printf("\t<%v> %s!\n", getTime(), str)
	fmt.Println("********************************************")
}

func getPdfFilePath() []string {
	var files []string
	root := "../../downloadsPDF/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Index(path, "pdf") != -1 {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

func getPdfId(path string) string {
	reg := regexp.MustCompile(`\d+`)
	id := reg.FindString(path)
	return id
}
