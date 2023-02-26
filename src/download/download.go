package download

import (
	"encoding/csv"
	"fmt"
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
		var csvRec CsvFile
		saveFileName := ""
		// 每条记录放入csvRec结构体
		if filepath == "../material/company-file-data/company-file-V3.csv" {
			csvRec = CsvFile{
				Id:             record[0],
				CompanyLink:    record[1],
				CompanyName:    record[2],
				FileAddr:       record[4],
				FileName:       record[9],
				DisclosureDate: record[10]}
			saveFileName = csvRec.Id + "-" + csvRec.FileName + ".pdf"
		} else {
			csvRec = CsvFile{
				Id:             record[0],
				CompanyLink:    "none",
				CompanyName:    record[1],
				FileAddr:       record[5],
				FileName:       record[4],
				DisclosureDate: record[6]}
			saveFileName = csvRec.FileName + ".pdf"
		}
		index := csvRec.Id
		fileUrl := csvRec.FileAddr
		saveFileBasePath := "../../downloadsPDF/"
		saveFilePath := saveFileBasePath + saveFileName
		downOneFile(fileUrl, saveFilePath, index)
		//
		// Todo: 必要内容保存到数据库
		// pf := database.PdfFile{
		// 	ID:          csvRec.Id,
		// 	CompanyLink: csvRec.CompanyLink,
		// 	CompanyName: csvRec.CompanyName,
		// 	FileName:    saveFileName,
		// }
		// database.Db.Create(&pf)
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
