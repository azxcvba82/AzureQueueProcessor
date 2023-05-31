package processor

import (
	"log"
	"time"
)

type AbstractProcessor struct {
	queueRequest QueueRequest
	logger       *QueueLogger
	queueName    string
}

func NewAbstractProcessor(queueRequest QueueRequest) *AbstractProcessor {
	newLogger := NewQueueLogger(queueRequest)
	queueName := queueRequest.RequestQueueName
	return &AbstractProcessor{queueRequest: queueRequest, logger: newLogger, queueName: queueName}
}

func (f *AbstractProcessor) Start(overrideProcess OverrideProcess) {
	defer func() {
		if e := recover(); e != nil {
			log.Fatal("Caught error: ", e)
		}
	}()
	f.PreProcessAction()
	overrideProcess.Process()
	f.PostProcessAction()
	f.LogAndCleanupAction()
}

func (f *AbstractProcessor) PreProcessAction() {
	f.logger.Log("Processor Processing Starting ..." + f.queueName)
	f.logger.Log("startTime UTC: " + time.Now().UTC().String())
}

func (f *AbstractProcessor) PostProcessAction() {
	f.logger.Log("endTime UTC: " + time.Now().UTC().String())
	f.logger.LogSave()
}

func (f *AbstractProcessor) LogAndCleanupAction() {
}

type OverrideProcess interface {
	Process()
}
