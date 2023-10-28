package main

import (
	"langxing.com/label-encryption/controllers"
	"log"
)

func main() {
	log.Println("启动显为加密程序...")

	_, err := controllers.TestSOAP()
	if err != nil {
		controllers.Quit(1)
	}
	code := controllers.RunBatchGeneratorController()
	controllers.Quit(code)
}
