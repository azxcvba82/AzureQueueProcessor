package utils

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/stretchr/objx"
)

func PostQueue(connString string, queueName string, message string) (res *HttpPost) {

	credential, err := NewSharedKeyCredential(connString)
	if err != nil {
		fmt.Println("NewSharedKeyCredential err")
	}

	URI := "https://" + credential.AccountName() + ".queue.core.windows.net/" + queueName + "/messages"
	template := "<QueueMessage><MessageText>" + message + "</MessageText></QueueMessage>"

	post := &HttpPost{
		URI:         URI,
		RequestBody: []byte(template),
	}
	err = credential.HttpPostRequest(post)
	post.ResponseBody = []byte(XML2JSON(string(post.ResponseBody)))
	// if err != nil {
	// 	fmt.Println("HttpPostRequest err")

	// } else {
	// 	fmt.Println(string(post.ResponseBody))

	// }
	return post
}

func GetQueue(connString string, queueName string) (res *HttpGet) {

	credential, err := NewSharedKeyCredential(connString)
	if err != nil {
		fmt.Println("NewSharedKeyCredential err")
	}

	URI := "https://" + credential.AccountName() + ".queue.core.windows.net/" + queueName + "/messages?numofmessages=1"

	get := &HttpGet{
		URI: URI,
	}
	err = credential.HttpGetRequest(get)
	get.ResponseBody = []byte(XML2JSON(string(get.ResponseBody)))
	// if err != nil {
	// 	fmt.Println("HttpGetRequest err")

	// } else {
	// 	fmt.Println(string(get.ResponseBody))

	// }
	return get
}

func PeekQueue(connString string, queueName string, params ...int) (res *HttpGet) {

	count := 32
	if len(params) > 0 {
		count = params[0]
	}

	credential, err := NewSharedKeyCredential(connString)
	if err != nil {
		fmt.Println("NewSharedKeyCredential err")
	}

	URI := "https://" + credential.AccountName() + ".queue.core.windows.net/" + queueName + "/messages?peekonly=true&numofmessages=" + strconv.Itoa(count)

	get := &HttpGet{
		URI: URI,
	}
	err = credential.HttpGetRequest(get)
	get.ResponseBody = []byte(XML2JSON(string(get.ResponseBody)))
	// if err != nil {
	// 	fmt.Println("HttpGetRequest err")

	// } else {
	// 	fmt.Println(string(get.ResponseBody))

	// }
	return get
}

func DeleteQueue(connString string, queueName string, messageid string, popreceipt string) (res *HttpGet) {

	credential, err := NewSharedKeyCredential(connString)
	if err != nil {
		fmt.Println("NewSharedKeyCredential err")
	}

	URI := "https://" + credential.AccountName() + ".queue.core.windows.net/" + queueName + "/messages/" + messageid + "?popreceipt=" + url.QueryEscape(popreceipt)

	delete := &HttpGet{
		URI: URI,
	}
	err = credential.HttpDeleteRequest(delete)
	delete.ResponseBody = []byte(XML2JSON(string(delete.ResponseBody)))
	// if err != nil {
	// 	fmt.Println("HttpGetRequest err")

	// } else {
	// 	fmt.Println(string(get.ResponseBody))

	// }
	return delete
}

func DeQueue(connString string, queueName string) (res *HttpGet) {
	get := GetQueue(connString, queueName)

	jobject, _ := objx.FromJSON(string(get.ResponseBody))

	if jobject.Get("QueueMessagesList").IsStr() == true {
		return nil
	}

	messageId := jobject.Get("QueueMessagesList.QueueMessage.MessageId").Str()
	popReceipt := jobject.Get("QueueMessagesList.QueueMessage.PopReceipt").Str()

	delete := DeleteQueue(connString, queueName, messageId, popReceipt)

	if delete.StatusCode != 204 {
		return delete
	}
	messageText := jobject.Get("QueueMessagesList.QueueMessage.MessageText").Str()
	get.ResponseBody = []byte(messageText)

	return get
}
