package parsepdf

import (
	"bufio"
	"fmt"
	"hello/database"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/ledongthuc/pdf"
)

var (
	totalLength, queLength, answerLength int
	indexTitle                           map[string]bool // 目录的各个标题
)

// pdf解析成txt
func ReadPdf(path string) {
	// {{{
	f, r, err := pdf.Open(path)
	if err != nil {
		log.Fatalln("[Parse.ReadPdf] Open pdf Failed! Err:", err)
	}
	defer f.Close()

	b, err := r.GetPlainText()

	f2, err := os.Create("../material/tmpFile/tmp.txt")
	if err != nil {
		log.Fatalln("[Parse.ReadPdf] Create tmp.txt failed! Err:", err)
	}
	defer f2.Close()

	_, err = io.Copy(f2, b)
	if err != nil {
		log.Fatalln("[Parse.ReadPdf] Copy pdf to tmp.txt Failed! Err:", err)
	}
	fmt.Println("[OK] Change Txt. File: ", path)
	// }}}
}

// 利用工具： `pdftotext` ,速度是自己写的6倍
func ReadPdfV2(path string) {
	//{{{
	textFilePath := "../material/tmpFile/tmp.txt"
	cmd := exec.Command("pdftotext", path, textFilePath)

	err := cmd.Run()
	if err != nil {
		log.Fatalln("[Parse.ReadPdfV2] change Failed! Err:", err)
	}
	fmt.Printf("[OK] Pdf Change to Txt. \nFileName: %v\n", path)
	//}}}
}

// 将一个pdf分割成两个txt临时文件，一个问询，一个回复
// 分词,分析
func DivideTwoParts() {
	// {{{
	// 1. 打开tmp文件准备操作
	txtFilePath := "../material/tmpFile/tmp.txt"
	f, err := os.Open(txtFilePath)
	if err != nil {
		log.Fatalln("[Parse.DivideTwoParts] open txtFile failed! Err:", err)
	}
	defer f.Close()

	// 2.创建两个临时文件，来写问询和回复的文本信息
	f1, err := os.Create("../material/tmpFile/tmpPart1.txt") // 问询
	defer f1.Close()
	f2, err := os.Create("../material/tmpFile/tmpPart2.txt") // 回复
	defer f2.Close()
	writeQue := bufio.NewWriter(f1)
	writeAns := bufio.NewWriter(f2)

	totalLength, queLength, answerLength = 0, 0, 0
	turn := 0 // 开关
	indexTitle = make(map[string]bool)

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("read err: ", err)
		}

		line = strings.TrimSpace(line)            // 去除\n等空白字符
		line = strings.Replace(line, " ", "", -1) // 去除空格

		// 将 是标题的行 加入标题map
		addTitleMap(line)

		// 主部分
		// 1.如果这行有标题,开始记问询数
		if indexTitle[line] || turn == 1 {
			for {
				line, err = reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				line = strings.TrimSpace(line)
				line = strings.Replace(line, " ", "", -1)
				if isAnswer(line) {
					answerLength += len([]rune(line))
					totalLength += len([]rune(line)) // 不管是什么都要加
					writeAns.WriteString(line + "\n")
					turn = -1
					break
				}
				queLength += len([]rune(line))
				totalLength += len([]rune(line)) // 不管是什么都要加
				writeQue.WriteString(line + "\n")
				writeAns.Flush()
				writeQue.Flush()
			}
		} else if isAnswer(line) || turn == -1 {
			// 2. 如果这行是回复的开始
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				line = strings.TrimSpace(line)
				line = strings.Replace(line, " ", "", -1)
				if indexTitle[line] {
					queLength += len([]rune(line))
					totalLength += len([]rune(line)) // 不管是什么都要加
					writeQue.WriteString(line + "\n")
					turn = 1
					break
				}
				answerLength += len([]rune(line))
				totalLength += len([]rune(line)) // 不管是什么都要加
				writeAns.WriteString(line + "\n")
				writeAns.Flush()
				writeQue.Flush()
			}
		}

		totalLength += len([]rune(line)) // 不管是什么都要加
	}
	fmt.Printf("分割成功：各文本数量如下")
	fmt.Printf("TotalLen:%v\nAnswerLength:%v\nQueLength:%v\n", totalLength, answerLength, queLength)
	fmt.Println(indexTitle)
	// }}}
}

// 结果部分写入数据库
func WriteSomeParseResToDB(id string) {
	pf := database.PdfFile{ID: id}
	count, text := countQueAmount()
	database.Db.Model(&pf).Updates(database.PdfFile{AllTextLen: totalLength,
		QueCount: count, QueText: text,
		QueryTextLen: queLength, AnswerTextLen: answerLength})
}

// ---------------------------------------------------
// 辅助函数

func init() {
}
func addTitleMap(line string) {
	//{{{
	bigTitle := map[string]bool{
		"一": true,
		"二": true,
		"三": true,
		"四": true,
		"五": true,
		"六": true,
		"七": true,
		"八": true,
		"九": true,
		"十": true,
	}
	titleDivideLine := "......."

	// 如果带有`....` ,不是大标题,不是小标题, 必须至少有一个数字,就加入map
	if strings.Index(line, titleDivideLine) != -1 &&
		(bigTitle[string([]rune(line)[0])] != true) &&
		isLittleTitle(line) &&
		isIncludeNum(line) {
		indexTitle[formatTitle(line)] = true // 将标题加入map
	}
	//}}}
}

// 是否string前段包含有数字
func isIncludeNum(s string) bool {
	for _, c := range formatTitle(s) {
		if '0' <= c && c <= '9' {
			return true
		}
	}
	return false
}

// 判断是不是小标题
func isLittleTitle(line string) bool {
	_, afterDotLine, ok := strings.Cut(line, ".")
	if !ok {
		return true
	}

	if '0' <= afterDotLine[0] && afterDotLine[0] <= '9' {
		return false
	}
	return true
}

// 提取.... 之前的东西并且去掉空格
func formatTitle(title string) string {
	res, _, _ := strings.Cut(title, "....")
	return res
}

// 回复 的开始标志
func isAnswer(line string) bool {
	// {{{
	if line == "回复：" || line == "回复:" || line == "【回复】" || line == "【发行人说明】" {
		return true
	}
	return false
	//}}}
}

// 计算问题的个数,返回文本 （计算的是目录中问题个数）
func countQueAmount() (int, string) {
	// {{{{
	count := 0
	text := ""
	for i := range indexTitle {
		count++
		text = text + "--" + i

	}
	return count, text
	//}}}
}
