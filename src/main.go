package main

import "hello/parsepdf"

func main() {
	// 下载
	// csvFilePath := "../material/company-file-data/company-file-all.csv"
	// download.ReadCsvAndDownLoad(csvFilePath)

	// 转换为pdf
	// pdf.DebugOn = true
	// content, err := parse.ReadPdf("../../downloadsPDF/三一重能股份有限公司-8-1 发行人及保荐机构回复意见.pdf") // Read local pdf file
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(content)

	// 分词

	parsepdf.Divide("../../txts/三生国健.txt")
	lawWordsFilePath := "../material/law-words-copy.txt"
	parsepdf.AddCompareWords(lawWordsFilePath)
	parsepdf.Count(parsepdf.Words)
}
