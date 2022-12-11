package download

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	companyNameMap map[string]bool
)

type csvFile struct {
	id             string
	companyLink    string
	companyName    string
	fileAddr       string
	fileName       string
	disclosureDate string
}

// 1.读取companyname, 将其放入切片中
func getCompany(filepath string) {
	companyNameMap = make(map[string]bool)
	// 打开文件
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("open %s failed\n : err%v", filepath, err)
		return
	}
	defer f.Close()

	// 按行读取
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("read err: ", err)
		}

		line = strings.TrimSpace(line) // 💫 readstring 会连带\n 保存,所以要剔除
		companyNameMap[line] = true    // 将每个公司置为true
	}
}

// 2.读取csv文件，返回需要的信息
func ReadCsvAndDownLoad(filepath string) {
	// 打开csv文件
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("open %s failed\n : err%v", filepath, err)
		os.Exit(1)
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
		csvRec := csvFile{record[0], record[1], record[2], record[3], record[4], record[5]}
		// 每条记录放入csvRec结构体
		if isCompany(csvRec.companyName, csvRec.fileName) {
			// 下载
			baseSaveUrl := "../../downloadsPDF/"
			fileUrl := csvRec.fileAddr

			// filename linux不能超过255字节
			saveFileName := baseSaveUrl + csvRec.companyName + "-" + csvRec.fileName[:] + ".pdf"
			if len(saveFileName) > 250 {
				saveFileName = baseSaveUrl + csvRec.companyName + "-" + csvRec.fileName[:100] + ".pdf"
			}

			downOneFile(fileUrl, saveFileName)
			// Todo: 必要内容保存到数据库

		}

	}

}

// 下载文件且重命名
func downOneFile(url string, saveFileName string) {
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
	fmt.Println("Start download...")
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("copy %s failed, err:%v\n", saveFileName, err)
	}
	fmt.Printf("[DOWNLOAD] %s successed!", saveFileName)
}

// 提前预存需要的公司名字
func init() {
	companyNameFilePath := "../material/companyName2.txt"
	getCompany(companyNameFilePath)

}

// 判断是否是需要的公司
func isCompany(companyName string, fileName string) bool {
	if companyNameMap[companyName] &&
		strings.Contains(fileName, "发行人") &&
		strings.Contains(fileName, "保荐机构") &&
		strings.Contains(fileName, "回复") {
		return true
	}

	return false
}
