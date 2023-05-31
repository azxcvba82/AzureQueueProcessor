package processor

import (
	"encoding/json"
	"main/utils"
)

type CurrencyConversionSync struct {
	*AbstractProcessor
}

func NewCurrencyConversionSyncProcessor(queueRequest QueueRequest) *CurrencyConversionSync {
	p := NewAbstractProcessor(queueRequest)
	return &CurrencyConversionSync{AbstractProcessor: p}
}

func (f *CurrencyConversionSync) Process() {
	f.logger.Log("processing ...  " + f.AbstractProcessor.queueName)
	s, _ := json.Marshal((f.AbstractProcessor.queueRequest))
	f.logger.Log("Request JObject" + string(s))

	// get conversion ratio
	URI := "https://openexchangerates.org/api/historical/2023-01-01.json?app_id=abba26dd7c8c40448b5006a312a8a411&base=USD"

	get := &utils.HttpGet{
		URI: URI,
	}
	err := utils.HttpGetRequest(get)
	if err != nil {
		f.logger.Log(err.Error())
	}
	f.logger.Log("raw data" + string(get.ResponseBody)[:250])

	// var model []interface{}
	// queryString := `SELECT * FROM Table `
	// err := utils.SQLQuery(&model, utils.GetSQLConnectString(), queryString)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
}
