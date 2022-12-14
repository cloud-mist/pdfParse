package download

import (
	"encoding/csv"
	"fmt"
	"hello/database"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	companyNameMap map[string]bool
)

type CsvFile struct {
	Id             string
	CompanyLink    string
	CompanyName    string
	FileAddr       string
	FileName       string
	DisclosureDate string
}

// 2.读取csv文件
func ReadCsvAndDownLoad(filepath string) {
	// {{{
	// 打开csv文件
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Open %s failed\n : err%v", filepath, err)
	}
	defer f.Close()

	// 读取,获得每条记录
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// 下载
		csvRec := CsvFile{
			Id:             record[0],
			CompanyLink:    record[1],
			CompanyName:    record[2],
			FileAddr:       record[4],
			FileName:       record[9],
			DisclosureDate: record[10]}
		// 每条记录放入csvRec结构体
		index := csvRec.Id
		saveFileName := index + "-" + csvRec.FileName + ".pdf"
		// fileUrl := csvRec.FileAddr
		// saveFileBasePath := "../../downloadsPDF/"
		// saveFilePath := saveFileBasePath + saveFileName
		// downOneFile(fileUrl, saveFilePath, index)
		// downOneFile(fileUrl, saveFileName, index) // 当前文件夹保存

		// Todo: 必要内容保存到数据库
		pf := database.PdfFile{
			ID:          csvRec.Id,
			CompanyLink: csvRec.CompanyLink,
			CompanyName: csvRec.CompanyName,
			FileName:    saveFileName,
		}
		database.Add2db(pf)
	}
	// }}}
}

// 下载文件且重命名
func downOneFile(url string, saveFileName string, index string) {
	// {{{
	// 请求
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("resp err : ", err)
	}
	defer resp.Body.Close()

	// 创建一个文件用于保存
	out, err := os.Create(saveFileName)
	if err != nil {
		log.Fatalln("Create File Err:", err)
	}
	defer out.Close()

	// 将响应copy到文件
	fmt.Printf("File<%s> Start Download...\n", index)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Copy <%s>%s Failed, err:%v\n", index, saveFileName, err)
	}
	fmt.Printf("[DOWNLOAD] %s Succeeded!\n", saveFileName)
	// }}}
}
