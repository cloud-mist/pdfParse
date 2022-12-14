package main

import "hello/download"

func main() {
	// 下载
	// csvFilePath := "company-file-V3.csv"
	csvFilePath := "../material/company-file-data/company-file-V3.csv"
	download.ReadCsvAndDownLoad(csvFilePath)

	// 分词
	// parsepdf.Divide("../../txts/三生国健.txt")
	// // parsepdf.FilterStopWords()
	// lawWordsFilePath := "../material/wordsFiles/law-words.txt"
	// parsepdf.AddCompareWords(lawWordsFilePath)
	// parsepdf.Count(parsepdf.PdfResWords)

}
