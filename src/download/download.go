package download

import (
	"fmt"
	"os"
)

var (
	companyNameMap map[string]bool
)

// 1.读取companyname, 将其放入切片中
func getCompany(filepath string) {

	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("open %s failed\n : err%v", filepath, err)
		return
	}

}

// 2.读取csv文件，返回需要的信息
func readCsv(filepath string) {
	f, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("open %s failed\n : err%v", filepath, err)
		return
	}
}

// 3. 对符合条件的公司问询函进行下载
func Download() {

}
