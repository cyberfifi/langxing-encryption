package controllers

import (
	"langxing.com/label-encryption/clients"
	"log"
	"os"
	"strings"
	"time"
)

func TestSOAP() (string, error) {
	res, err := clients.GetHelloResult()
	if err != nil {
		log.Printf("SOAP客户端测试出现错误，退出程序: %v", err)
		return "", err
	}
	log.Printf("SOAP客户端测试成功。结果：%s", res)
	return res, nil
}

func GetResult(input string) (string, string, error) {
	input = strings.Trim(input, " ")
	//if len(input) <= 20 {
	//	return "", "", errors.New("输入无效。请核查你输入的信息是否正确")
	//}
	res, err := clients.EncryptText(input)
	if err != nil {
		log.Printf("加密失败: %v \n", err.Error())
		return "", "", err
	}
	//verifyCode, err := clients.GetVerifyCode(input[0:18])
	//if err != nil {
	//	log.Printf("获取校验位失败: %v \n", err.Error())
	//	return res, "", err
	//}
	//return res, verifyCode, nil
	return res, "", nil
}

func Quit(code int) {
	log.Println("正在退出程序...")
	time.Sleep(3 * time.Second)
	os.Exit(code)
}
