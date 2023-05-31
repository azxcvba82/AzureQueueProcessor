package main

import (
	"encoding/json"
	"fmt"
	"main/processor"
	"main/utils"
	"os"
	"time"
)

func main() {
	//TestPost()

	connString := GetConn()

	for {
		time.Sleep(time.Second * 5)

		go func() {
			res := utils.DeQueue(connString, "demo1")
			if res != nil {

				var req processor.QueueRequest
				if err := json.Unmarshal([]byte(res.ResponseBody), &req); err != nil {
					fmt.Println("err")
					return
				}
				go func() {
					p := processor.NewWeatherforecastProcessor(req)
					p.Start(p)
				}()
			}
		}()

	}

}

func GetConn() string {
	return os.Getenv("STORAGE_CONNECTION_STRING")
}

func TestPost() {
	connString := GetConn()
	data := "{}"
	res := utils.PostQueue(connString, "demo1", data)
	fmt.Println(string(res.ResponseBody))
}
