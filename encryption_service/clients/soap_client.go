package clients

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	URL = "http://localhost:8090/axis2/services/EncryptWS"
)

type SoapHeader struct {
	XMLName xml.Name `xml:"soapenv:Header"`
}

// Encrypt endpoint

type SoapEncryptRequestEnvelop struct {
	XMLName   xml.Name `xml:"soapenv:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soapenv,attr"`
	XMLNsWser string   `xml:"xmlns:wser,attr"`
	Header    SoapHeader
	Body      EncryptRequestBody
}

type EncryptRequestBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Encrypt EncryptRequest
}

type EncryptRequest struct {
	XMLName xml.Name `xml:"wser:encrypt"`
	Data    string   `xml:"wser:data"`
}

type SoapEncryptResponse struct {
	XMLName xml.Name
	Body    SoapEncryptResponseBody
}

type SoapEncryptResponseBody struct {
	XMLName         xml.Name
	EncryptResponse EncryptResponse `xml:"encryptResponse"`
}

type EncryptResponse struct {
	XMLName xml.Name `xml:"encryptResponse"`
	Return  string   `xml:"return"`
}

// Hello endpoint

type SoapHelloRequestEnvelop struct {
	XMLName   xml.Name `xml:"soapenv:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soapenv,attr"`
	XMLNsWser string   `xml:"xmlns:wser,attr"`
	Header    SoapHeader
	Body      HelloRequestBody
}

type HelloRequestBody struct {
	XMLName xml.Name `xml:"soapenv:Body"`
	Hello   HelloRequest
}

type HelloRequest struct {
	XMLName xml.Name `xml:"wser:hello"`
	Name    string   `xml:"wser:name"`
}

type SoapHelloResponse struct {
	XMLName xml.Name
	Body    SoapHelloResponseBody
}

type SoapHelloResponseBody struct {
	XMLName       xml.Name
	HelloResponse HelloResponse `xml:"helloResponse"`
}

type HelloResponse struct {
	XMLName xml.Name `xml:"helloResponse"`
	Return  string   `xml:"return"`
}

// GetVerifyCode endpoint

type SoapGetVerifyCodeRequestEnvelop struct {
	XMLName   xml.Name `xml:"soapenv:Envelope"`
	XMLNsSoap string   `xml:"xmlns:soapenv,attr"`
	XMLNsWser string   `xml:"xmlns:wser,attr"`
	Header    SoapHeader
	Body      GetVerifyCodeRequestBody
}

type GetVerifyCodeRequestBody struct {
	XMLName       xml.Name `xml:"soapenv:Body"`
	GetVerifyCode GetVerifyCodeRequest
}

type GetVerifyCodeRequest struct {
	XMLName xml.Name `xml:"wser:getVerifyCode"`
	Data    string   `xml:"wser:data"`
}

type SoapGetVerifyCodeResponse struct {
	XMLName xml.Name
	Body    SoapGetVerifyCodeResponseBody
}

type SoapGetVerifyCodeResponseBody struct {
	XMLName               xml.Name
	GetVerifyCodeResponse GetVerifyCodeResponse `xml:"getVerifyCodeResponse"`
}

type GetVerifyCodeResponse struct {
	XMLName xml.Name `xml:"getVerifyCodeResponse"`
	Return  string   `xml:"return"`
}

func GetHelloResult() (string, error) {
	log.Println("开始测试SOAP Hello 服务器...")
	helloRq := SoapHelloRequestEnvelop{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsWser: "http://wserver",
		Body: HelloRequestBody{
			Hello: HelloRequest{
				Name: "hello",
			},
		},
	}
	payload, err := xml.MarshalIndent(helloRq, "", "  ")
	fmt.Println("SOAP请求----")
	fmt.Println(string(payload))
	timeout := 30 * time.Second
	client := http.Client{
		Timeout: timeout,
	}
	ws := "http://localhost:8090/axis2/services/EncryptWS"
	req, err := http.NewRequest("POST", ws, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", "POST")
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("SOAP结果:")
	fmt.Println(string(bodyBytes))
	defer response.Body.Close()
	var result SoapHelloResponse
	err = xml.Unmarshal(bodyBytes, &result)
	if err != nil {
		//log.Printf("测试出现错误，退出程序 %v", err.Error())
		//time.Sleep(3 * time.Second)
		//os.Exit(1)
		return "", err
	}

	fmt.Println("处理后的测试结果:")
	fmt.Println(result.Body.HelloResponse.Return)
	log.Println("完成SOAP服务器测试...")
	return result.Body.HelloResponse.Return, nil
}

func GetVerifyCode(text string) (string, error) {
	log.Println("开始获取校验位...")
	envelop := SoapGetVerifyCodeRequestEnvelop{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsWser: "http://wserver",
		Body: GetVerifyCodeRequestBody{
			GetVerifyCode: GetVerifyCodeRequest{
				Data: text,
			},
		},
	}
	payload, err := xml.MarshalIndent(envelop, "", "  ")
	fmt.Println("SOAP请求:")
	fmt.Println(string(payload))
	fmt.Println("---------")

	timeout := 30 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", "POST")
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("SOAP结果:")
	fmt.Println(string(bodyBytes))
	defer response.Body.Close()
	var result SoapGetVerifyCodeResponse
	err = xml.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", err
	}
	return result.Body.GetVerifyCodeResponse.Return, nil
}

func EncryptText(text string) (string, error) {
	envelop := SoapEncryptRequestEnvelop{
		XMLNsSoap: "http://schemas.xmlsoap.org/soap/envelope/",
		XMLNsWser: "http://wserver",
		Body: EncryptRequestBody{
			Encrypt: EncryptRequest{
				Data: text,
			},
		},
	}
	payload, err := xml.MarshalIndent(envelop, "", "  ")
	fmt.Println("SOAP请求:")
	fmt.Println(string(payload))
	fmt.Println("----------------------")

	timeout := 30 * time.Second
	client := http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", "POST")
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("SOAP结果:")
	fmt.Println(string(bodyBytes))
	fmt.Println("----------------------")
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)
	var result SoapEncryptResponse
	err = xml.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", err
	}
	return result.Body.EncryptResponse.Return, nil
}
