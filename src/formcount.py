import pdfplumber

filePath = "/home/shawn/study/Project/crawl/downloadsPDF/77-688068.SH_热景生物_四轮反馈回复.pdf"
# filePath = "/home/shawn/Desktop/myPhoto/学生基本信息表zx.pdf"


def count(filePath):
    oneLineHeight, imgHeight, tableRowRes = 0, 0, 0
    totalRow = 0
    manyImg = []
    with pdfplumber.open(filePath) as pdf:
        pages = pdf.pages
        for page in pages:
            # image
            images = page.images
            for image in images:
                for k, v in image.items():
                    if k == "height":
                        manyImg.append(v)

            # table
            tables = page.find_tables()
            for table in tables:
                tableRowRes += len(table.rows)

            # totalRow
            texts = page.extract_text()
            for text in texts:
                if text == "\n":
                    totalRow += 1

            # oneLine
            chars = page.chars
            for char in chars:
                for k, v in char.items():
                    if k == "height":
                        oneLineHeight = v
                        break
                break
    # 计算imgHeight
    manyImg = set(manyImg)
    for i in manyImg:
        imgHeight += i

    print(
        f"总行数：{totalRow}, 表格行数：{tableRowRes}, 图片行数： {imgHeight/oneLineHeight:.2f} \n图片高度：{imgHeight:.2f},一行的高度：{oneLineHeight:.2f}"
    )


count(filePath)
