package parsepdf

import (
	"bufio"
	"fmt"
	"hello/database"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
)

var (
	totalLength, queLength, answerLength int
	indexTitle                           map[string]bool // 目录的各个标题
	signTitle                            int
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
func ReadPdfV2(path string, textFilePath string) {
	//{{{
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

	// 全局变量初始化
	totalLength, queLength, answerLength = 0, 0, 0
	turn := 0 // 开关
	indexTitle = make(map[string]bool)
	signTitle = 0

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

func addTitleMapV2(line string) {
	//{{{
	reg := regexp.MustCompile(`(问题)?\d+.*`)
	res := reg.FindString(line)

	titleDivideLine := "......."
	if strings.Index(line, titleDivideLine) != -1 && strings.Index(line, "目录") == -1 {
		if len(res) != 0 {
			indexTitle[formatTitle(line)] = true
		}
	}
	//}}}
}

func addTitleMap(line string) {
	titleDivideLine := "......."

	// 只有在是目录标题，且不是目录两字的标题的时候运行
	if strings.Index(line, titleDivideLine) != -1 &&
		strings.Index(line, "目录") == -1 && strings.Index(line, "目 录") == -1 {
		indexTitleLength := len(indexTitle)

		if indexTitleLength == 0 {
			indexTitle[formatTitle(line)] = true
		} else if (indexTitleLength == 1) && (signTitle == 0) {
			if hasTwoLevelTitle(line) {
				if hasPrombleInTitle(line) {
					// 如果 有两级标题，且问题在二级标题里，那么清空一级标题，且标记;
					indexTitle = make(map[string]bool)
					indexTitle[formatTitle(line)] = true
					signTitle = 22
				} else {
					//                    问题在一级标题，那么标记
					signTitle = 21
				}
			} else {
				indexTitle[formatTitle(line)] = true
				signTitle = 1
			}
		} else if indexTitleLength >= 1 && (signTitle != 0) {
			// 已经是第三个标题及之后
			// 如果是只有一级标题，这个标题和前面的标题类型一致，那么就添加

			if !hasTwoLevelTitle(line) {
				indexTitle[formatTitle(line)] = true
			}
		}
	}
}

// Helper: addTitleMap
func hasPrombleInTitle(line string) bool {
	problem := "问题"
	if strings.Index(line, problem) == -1 {
		return false
	}
	return true
}

// Helper: addTitleMap
// 只有在长度为1时，计算是否有两级标题
func hasTwoLevelTitle(secTitle string) bool {
	tmpMap1 := map[string]bool{
		"一": true, "二": true, "三": true, "四": true, "五": true,
		"零": true, "九": true, "八": true, "七": true, "六": true,
	}
	tmpMap2 := map[string]bool{
		"1": true, "2": true, "3": true, "4": true, "5": true,
		"9": true, "0": true, "8": true, "7": true, "6": true,
	}
	firstTitle := ""
	for i := range indexTitle {
		firstTitle = i
		break
	}
	firstTitlef, secTitlef := string([]rune(firstTitle)[0]), string([]rune(secTitle)[0]) // 获取前两个标题的第一个字
	// 如果前两个标题 第一个字相同或者格式相同, 那么就只有一级标题，否则是两级
	if (firstTitlef == secTitlef) ||
		(tmpMap1[firstTitlef] && tmpMap1[secTitlef]) ||
		(tmpMap2[firstTitlef] && tmpMap2[secTitlef]) {

		return false
	}
	return true
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
// TODO: 未实现 （1）单纯一个回复 （2）直接回答
func isAnswer(line string) bool {
	// {{{
	AnswerSign := map[string]bool{
		"回复hf:": true, "回复hf：": true,
	}
	if AnswerSign[line] {
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
