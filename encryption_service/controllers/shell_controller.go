package controllers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func RunShell() int {
	for {
		log.Println("============================================")
		input := StringPrompt("请输入需要加密的字符串 (如要退出程序，请输入\"exit\"): ")
		if input == "exit" {
			return 0
		}
		encryptedText, verifyCode, err := GetResult(input)
		log.Println("============================================")
		if err != nil {
			log.Printf("系统错误：%v", err)
		}
		if encryptedText != "" {
			log.Printf("加密结果: %s\n", encryptedText)
		}
		if verifyCode != "" {
			log.Printf("校验位: %s\n", verifyCode)
		}
	}
}

// StringPrompt asks for a string value using the label
func StringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		_, err := fmt.Fprint(os.Stderr, label+" ")
		if err != nil {
			return ""
		}
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}
