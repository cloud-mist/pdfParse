package parsepdf

import (
	"bufio"
	"fmt"
	"hello/database"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ego/gse"
)

// 1.加完词库
// 2. 再分词
var (
	pdfResWords     []string        // pdf的分词过滤结果
	compareWordsMap map[string]bool // 当前需要放入的词库
	total           map[string]int  // 统计的词汇情况
	vocLen          int
	newJieba        gse.Segmenter
)

func init() {
	lawWordsFilePath := "../material/wordsFiles/law-words.txt"
	accountingWordsFilePath := "../material/wordsFiles/accounting-words.txt"
	financialWordsFilePath := "../material/wordsFiles/financial-words.txt"
	stopWordsPath := "../material/wordsFiles/stop-words-copy.txt"

	newJieba, _ = gse.New()
	newJieba.LoadDict(lawWordsFilePath + "," + accountingWordsFilePath + "," + financialWordsFilePath)
	// newJieba.LoadDict(OtherWordsDic()) // 其他词库
	newJieba.LoadStop(stopWordsPath)
}

func Divide(txtFilePath string) {
	// 分词
	f, err := os.Open(txtFilePath)
	if err != nil {
		log.Fatalf("[Parse.Divide] Open TXTFile Failed! err: %v\n", err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	pdfResWords = []string{} // 初始化
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err: ", err)
		}
		for _, v := range newJieba.Trim(newJieba.Cut(line, true)) {
			pdfResWords = append(pdfResWords, v)
			// words []string ， true hmm开启，
		}
	}
}

// 匹配计数

func Count() {
	total = make(map[string]int)

	// 暴力匹配
	vocLen = 0
	for i := range pdfResWords {
		if compareWordsMap[pdfResWords[i]] {
			vocLen++
			total[pdfResWords[i]]++
		}
	}
	fmt.Println("匹配的单词总数：", vocLen)
	fmt.Println("匹配的单词情况：")
	// table := tablewriter.NewWriter(os.Stdout)
	// for i, v := range total {
	// 	table.Append([]string{i, strconv.Itoa(v)})
	// }
	// table.Render()
}

func AddCompareWords(wordsFilePath string) {
	//{{{
	compareWordsMap = make(map[string]bool)
	f, err := os.Open(wordsFilePath)
	if err != nil {
		log.Fatalf("[Parse.AddCompareWords] Open wordsFile Failed! err: %v\n", err)
	}
	defer f.Close()

	// 读取每个词
	reader := bufio.NewReader(f)
	for {
		tmp, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("read err: ", err)
		}
		word := strings.Fields(tmp)
		compareWordsMap[word[0]] = true
	}
	//}}}
}

// 其他字典路径拼接
func OtherWordsDic() string {
	//{{{
	var files string
	root := "../../cacl2/dicts/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path[len(path)-3:] == "txt" { // 如果后缀是txt,才加入
			files = files + "," + path
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
	//}}}
}

// 插入数据库
func WriteWordsVocNum(id, sign string) {
	pf := database.PdfFile{ID: id}
	if sign == "law" {
		database.Db.Model(&pf).Updates(database.PdfFile{LawVocLen: vocLen})
	} else if sign == "account" {
		database.Db.Model(&pf).Updates(database.PdfFile{AccountVocLen: vocLen})
	}
}
