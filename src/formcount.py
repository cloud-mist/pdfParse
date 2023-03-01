import pdfplumber
import os

allPdfPath = "/home/shawn/study/Project/crawl/downloadsPDF/"


def count(filePath):
    oneLineHeight, imgHeight, tableRowRes = 0, 0, 0
    totalRow = 0
    manyImg = []
    with pdfplumber.open(filePath) as pdf:
        pages = pdf.pages
        length = len(pages)
        for page in pages[1 : length - 5]:
            # image
            images = page.images
            for image in images:
                for k, v in image.items():
                    if k == "height":
                        # print(f"img: page:{page}, height:{v}")
                        manyImg.append(v)

        for page in pages[1:]:
            # table
            tables = page.find_tables()
            for table in tables:
                # print(f"table: page:{page}")
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
    # print(manyImg)
    for i in manyImg:
        imgHeight += i

    # print(
    #     f"总行数：{totalRow}, 表格行数：{tableRowRes}, 图片行数： {imgHeight/oneLineHeight:.2f} \n图片高度：{imgHeight:.2f},一行的高度：{oneLineHeight:.2f}"
    # )
    return [totalRow, tableRowRes, imgHeight, oneLineHeight]


# ----------------------------------------------------
for fileName in os.listdir(allPdfPath):
    filePath = allPdfPath + fileName
    print(fileName)
    # print(filePath)
    # count(filePath)
    totalRow, tableRowRes, imgHeight, oneLineHeight = count(filePath)
    with open("./res.csv", "a", encoding="utf-8") as f:
        f.write(
            f"{fileName},{totalRow},{tableRowRes},{imgHeight/oneLineHeight:.2f},{imgHeight:.2f},{oneLineHeight:.2f}\n"
        )
