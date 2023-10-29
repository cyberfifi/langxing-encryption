package controllers

import (
	"errors"
	"fmt"
	"langxing.com/label-encryption/clients"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	VendorCode     = "LFE86"
	BanCode        = "A1"
	SeparationCode = "*"
	FillCode       = "001"
)

type RunOneBatchRequest struct {
	BatchIndex int
	FileName   string
	RawDate    string
	StartId    string
	PartCode   string
}

type RunOneBatchResponse struct {
	PreviousDate     string
	NextId           string
	PreviousPartCode string
	ShouldContinue   bool
}

type Result struct {
	VendorCode     string
	DateCode       string
	BanCode        string
	IdCode         string
	PartCode       string
	SeparationCode string
	FillCode       string
	VerifyCode     string
	EncryptedCode  string
}

func RunBatchGeneratorController() int {
	fileName := StringPrompt("请输入文件名 (例子 001): ")
	if len(fileName) == 0 {
		log.Print("输入的文件名无效")
		return 1
	}
	if !strings.HasSuffix(fileName, ".txt") {
		fileName = fileName + ".txt"
	}
	curBatchIndex := 0
	initRequest := RunOneBatchRequest{
		BatchIndex: curBatchIndex,
		FileName: fileName,
	}
	res, err := RunOneBatch(initRequest)
	if err != nil {
		return 1
	}
	curBatchIndex += 1
	for res.ShouldContinue {
		req := RunOneBatchRequest{
			BatchIndex: curBatchIndex,
			FileName: fileName,
			StartId:  res.NextId,
			PartCode: res.PreviousPartCode,
			RawDate:  res.PreviousDate,
		}
		res, err = RunOneBatch(req)
		if err != nil {
			return 1
		}
		curBatchIndex += 1
	}
	return 0
}

func RunOneBatch(req RunOneBatchRequest) (*RunOneBatchResponse, error) {
	res := RunOneBatchResponse{}
	partCodeText := "请输入零件码 (例子 8101010-M01-C00/C)"
	if req.PartCode != "" {
		partCodeText = fmt.Sprintf("请输入零件码 (上次输入 %s): ", req.PartCode)
	}
	partCode := StringPrompt(partCodeText)
	if len(partCode) == 0 {
		log.Print("输入的零件码无效")
		return nil, errors.New("输入的零件码无效")
	}

	rawDateText := "请输入日期 (例子 20231027): "
	if req.RawDate != "" {
		rawDateText = fmt.Sprintf("请输入日期 (上次输入 %s): ", req.RawDate)
	}
	rawDate := StringPrompt(rawDateText)
	if len(rawDate) != 8 {
		log.Print("输入的日期无效")
		return nil, errors.New("输入的日期无效")
	}
	dateInput, err := time.Parse("20060102", rawDate)
	if err != nil {
		log.Print("输入的日期无效")
		return nil, errors.New("输入的日期无效")
	}

	rawStartIdText := "请输入起始流水号 (例子 0000003): "
	if req.StartId != "" {
		rawStartIdText = fmt.Sprintf("请输入起始流水号 (上次结尾流水号 %s): ", req.StartId)
	}
	rawStartId := StringPrompt(rawStartIdText)
	if len(rawStartId) > 7 {
		log.Print("输入的起始流水号无效")
		return nil, errors.New("输入的起始流水号无效")
	}
	startId, err := strconv.Atoi(rawStartId)
	if err != nil {
		log.Print("输入的起始流水号无效")
		return nil, errors.New("输入的起始流水号无效")
	}

	rawTotalNum := StringPrompt("请输入所需标签数量 (例子 25): ")
	totalNum, err := strconv.Atoi(rawTotalNum)
	if err != nil {
		log.Print("输入的标签数量无效")
		return nil, errors.New("输入的标签数量无效")
	}

	f, err := os.OpenFile(req.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return nil, errors.New("无法打开文件")
	}
	defer f.Close()

	i := 0
	curId := startId
	if req.BatchIndex > 0 {
		if _, err := f.WriteString("\n"); err != nil {
			log.Println(err)
			return nil, errors.New("无法将起始行写入文件")
		}
	}
	for i < totalNum {
		dateCode := GetDateCode(dateInput)
		idCode := fmt.Sprintf("%07d", curId)

		codeToEncrypt := fmt.Sprintf("%s%s%s%s%s%s%s\n", VendorCode, dateCode, BanCode, idCode, partCode,
			SeparationCode, FillCode)
		encryptedCode, err := clients.EncryptText(codeToEncrypt)
		if err != nil {
			log.Printf("无法加密文字: %s. %v", codeToEncrypt, err)
			return nil, errors.New("无法加密文字")
		}

		codeToGetVerifyCode := fmt.Sprintf("%s%s%s%s", VendorCode, dateCode, BanCode, idCode)
		verifyCode, err := clients.GetVerifyCode(codeToGetVerifyCode)
		if err != nil {
			log.Printf("无法获取校验位: %s. %v", codeToGetVerifyCode, err)
			return nil, errors.New("无法获取校验位")
		}
		res := Result{
			VendorCode:     VendorCode,
			DateCode:       dateCode,
			BanCode:        BanCode,
			IdCode:         idCode,
			VerifyCode:     verifyCode,
			PartCode:       partCode,
			SeparationCode: SeparationCode,
			FillCode:       FillCode,
			EncryptedCode:  encryptedCode,
		}
		log.Println(res)
		resList := []string{
			VendorCode, dateCode, BanCode, idCode, verifyCode,
			partCode, SeparationCode, FillCode, encryptedCode}

		resString := strings.Join(resList, "\t")
		stringToWrite := resString + "\n"
		if i == totalNum-1 {
			stringToWrite = resString
		}
		if _, err := f.WriteString(stringToWrite); err != nil {
			log.Println(err)
			return nil, errors.New("无法将结果写入文件")
		}

		curId += 1
		i += 1
	}

	shouldContinueText := fmt.Sprintf(
		"已经成功生成一批记录，共计%d行。是否继续输入更多结果到同一文件？(请输入YES或者NO)",
		totalNum)
	shouldContinue := strings.ToUpper(StringPrompt(shouldContinueText))
	res.NextId = fmt.Sprintf("%07d", curId-1)
	res.PreviousDate = rawDate
	res.PreviousPartCode = partCode
	if shouldContinue == "YES" {
		log.Println("收到继续指令，准备下一批次输入...")
		res.ShouldContinue = true
	} else {
		log.Println("未收到继续指令，准备退出程序...")
		res.ShouldContinue = false
	}
	return &res, nil
}

func GetDateCode(date time.Time) string {
	year := date.Year()
	month := int(date.Month())
	day := date.Day()
	return fmt.Sprintf("%s%s%s", GetYearCode(year), GetMonthCode(month), GetDayCode(day))
}

func GetYearCode(num int) string {
	yearMap := map[int]string{
		2023: "P",
		2024: "R",
		2025: "S",
	}
	return yearMap[num]
}

func GetMonthCode(num int) string {
	monthMap := map[int]string{
		1:  "1",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "A",
		11: "B",
		12: "C",
	}
	return monthMap[num]
}

func GetDayCode(num int) string {
	code := fmt.Sprintf("%02d", num)
	return code
}
