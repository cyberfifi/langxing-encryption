package controllers

import (
	"encoding/json"
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
	PartCode       = "U8101010-M21-C00/C"
	SeparationCode = "*"
	FillCode       = "001"
)

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
	rawDate := StringPrompt("请输入日期 (例子 20231027): ")
	if len(rawDate) != 8 {
		log.Print("输入的日期无效")
		return 1
	}
	dateInput, err := time.Parse("20060102", rawDate)
	if err != nil {
		log.Print("输入的日期无效")
		return 1
	}
	rawStartId := StringPrompt("请输入起始流水号 (例子 0000003): ")
	if len(rawStartId) != 7 {
		log.Print("输入的起始流水号无效")
		return 1
	}
	startId, err := strconv.Atoi(rawStartId)
	if err != nil {
		log.Print("输入的起始流水号无效")
		return 1
	}
	rawTotalNum := StringPrompt("请输入所需标签数量 (例子 25): ")
	totalNum, err := strconv.Atoi(rawTotalNum)
	if err != nil {
		log.Print("输入的标签数量无效")
		return 1
	}

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return 1
	}
	defer f.Close()

	i := 0
	curId := startId
	for i < totalNum {
		dateCode := GetDateCode(dateInput)
		idCode := fmt.Sprintf("%07d", curId)

		codeToEncrypt := fmt.Sprintf("%s%s%s%s%s%s%s\n", VendorCode, dateCode, BanCode, idCode, PartCode,
			SeparationCode, FillCode)
		encryptedCode, err := clients.EncryptText(codeToEncrypt)
		if err != nil {
			log.Printf("无法加密文字: %s. %v", codeToEncrypt, err)
			return 1
		}

		codeToGetVerifyCode := fmt.Sprintf("%s%s%s%s", VendorCode, dateCode, BanCode, idCode)
		verifyCode, err := clients.GetVerifyCode(codeToGetVerifyCode)
		if err != nil {
			log.Printf("无法获取校验位: %s. %v", codeToGetVerifyCode, err)
			return 1
		}

		res := Result{
			VendorCode:     VendorCode,
			DateCode:       dateCode,
			BanCode:        BanCode,
			IdCode:         idCode,
			VerifyCode:     verifyCode,
			PartCode:       PartCode,
			SeparationCode: SeparationCode,
			FillCode:       FillCode,
			EncryptedCode:  encryptedCode,
		}

		resJson, err := json.Marshal(res)
		if err != nil {
			fmt.Println("JSON处理出错")
		}

		textToAppend := string(resJson)

		if _, err := f.WriteString(textToAppend + "\n"); err != nil {
			log.Println(err)
			return 1
		}

		curId += 1
		i += 1
	}
	return 0
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
