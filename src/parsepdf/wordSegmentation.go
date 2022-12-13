package parsepdf

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-ego/gse"
	"github.com/olekukonko/tablewriter"
	"github.com/yanyiwu/gojieba"
)

// 1.加完词库
// 2. 再分词
var (
	PdfResWords []string // pdf的分词结果
	// FilterResWords  []string        // 分词过滤后的结果
	compareWordsMap map[string]bool // 当前需要放入的词库
	total           map[string]int  // 统计的词汇情况
)

func Divide(txtFilePath string) {
	// newJieba := gojieba.NewJieba()
	// defer newJieba.Free()
	lawWordsFilePath := "../material/wordsFiles/law-words.txt"
	accountingWordsFilePath := "../material/wordsFiles/accounting-words.txt"
	financialWordsFilePath := "../material/wordsFiles/financial-words.txt"
	stopWordsPath := "../material/wordsFiles/stop-words-copy.txt"

	newJieba, _ := gse.New()
	newJieba.LoadDict(lawWordsFilePath + "," + accountingWordsFilePath + "," + financialWordsFilePath)
	newJieba.LoadStop(stopWordsPath)

	// 添加法律,会计，金融词库
	// lawWordsFilePath := "../material/wordsFiles/law-words.txt"
	// accountingWordsFilePath := "../material/wordsFiles/accounting-words.txt"
	// financialWordsFilePath := "../material/wordsFiles/financial-words.txt"
	// addWordsToDic(lawWordsFilePath, newJieba)
	// addWordsToDic(accountingWordsFilePath, newJieba)
	// addWordsToDic(financialWordsFilePath, newJieba)

	// 分词
	f, err := os.Open(txtFilePath)
	if err != nil {
		log.Fatalf("[Parse] Open TXTFile Failed! err: %v\n", err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err: ", err)
		}
		PdfResWords = newJieba.Trim(newJieba.Cut(line, true)) // words []string ， true hmm开启，
		fmt.Println(PdfResWords)
	}
}

// 添加词库
func addWordsToDic(wordsFilePath string, newJieba *gojieba.Jieba) {
	// 打开文件
	f, err := os.Open(wordsFilePath)
	if err != nil {
		log.Fatalf("[Parse] Open wordsFile Failed! err: %v\n", err)
	}
	defer f.Close()

	// 读取每个词
	reader := bufio.NewReader(f)
	for {
		word, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("read err: ", err)
		}

		word = strings.TrimSpace(word)
		newJieba.AddWord(word) // 添加词
	}
}

// 匹配计数
func Count(strSlice []string) {
	total = make(map[string]int)

	// 暴力匹配
	count := 0
	for i := range strSlice {
		if compareWordsMap[strSlice[i]] {
			count++
			total[strSlice[i]]++
		}
	}
	fmt.Println(count)
	table := tablewriter.NewWriter(os.Stdout)
	for i, v := range total {
		table.Append([]string{i, strconv.Itoa(v)})
	}
	table.Render()
}

func AddCompareWords(wordsFilePath string) {
	compareWordsMap = make(map[string]bool)
	f, err := os.Open(wordsFilePath)
	if err != nil {
		log.Fatalf("[Parse] Open wordsFile Failed! err: %v\n", err)
	}
	defer f.Close()

	// 读取每个词
	reader := bufio.NewReader(f)
	for {
		word, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("read err: ", err)
		}

		word = strings.TrimSpace(word)
		compareWordsMap[word] = true
	}
}

// 过滤stopwords
// func FilterStopWords() {
// 	stopWordsMap := make(map[string]bool)
// 	stopWordsPath := "../material/wordsFiles/stop-words-copy.txt"
// 	f, err := os.Open(stopWordsPath)
// 	if err != nil {
// 		log.Fatalf("[Filter] Open stopWordsFile Failed! err: %v\n", err)
// 	}
// 	defer f.Close()
//
// 	// 读取每个词
// 	reader := bufio.NewReader(f)
// 	for {
// 		word, err := reader.ReadString('\n')
// 		if err == io.EOF {
// 			break
// 		}
//
// 		if err != nil {
// 			fmt.Println("read err: ", err)
// 		}
//
// 		word = strings.TrimSpace(word)
// 		stopWordsMap[word] = true
// 	}
//
// 	// 过滤
// 	for i := range pdfResWords {
// 		if !stopWordsMap[pdfResWords[i]] {
// 			FilterResWords = append(FilterResWords, pdfResWords[i])
// 		}
// 	}
// }
