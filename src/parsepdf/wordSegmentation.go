package parsepdf

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/yanyiwu/gojieba"
)

// 1.加完词库
// 2. 再分词
var (
	Words           []string        // pdf的分词结果
	compareWordsMap map[string]bool // 当前需要放入的词库
	total           map[string]int  // 统计的词汇情况
)

func Divide(txtFilePath string) {
	newJieba := gojieba.NewJieba()
	defer newJieba.Free()

	// 添加法律,会计，金融词库
	lawWordsFilePath := "../material/law-words.txt"
	accountingWordsFilePath := "../material/accounting-words.txt"
	financialWordsFilePath := "../material/financial-words.txt"
	addWordsToDic(lawWordsFilePath, newJieba)
	addWordsToDic(accountingWordsFilePath, newJieba)
	addWordsToDic(financialWordsFilePath, newJieba)

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
		Words = newJieba.Cut(line, true) // words []string ， true hmm开启，
		fmt.Println(Words)
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
