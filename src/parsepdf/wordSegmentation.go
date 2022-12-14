package parsepdf

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-ego/gse"
	"github.com/olekukonko/tablewriter"
)

// 1.加完词库
// 2. 再分词
var (
	PdfResWords     []string        // pdf的分词过滤结果
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
	newJieba.LoadDict(OtherWordsDic()) // 其他词库
	newJieba.LoadStop(stopWordsPath)
	newJieba.AddToken("科创板", 10)

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
}

func OtherWordsDic() string {
	var files string
	root := "../../cacl2/dicts/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if path[len(path)-3:] == "txt" {
			files = files + "," + path
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}
