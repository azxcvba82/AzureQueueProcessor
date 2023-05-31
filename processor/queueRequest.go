package processor

type QueueRequest struct {
	RequestStorageConnectionString string
	RequestQueueName               string

	LogStorageConnectionString string
	LogContainerName           string
	LogFileName                string

	RequestTime string

	DBConnectionStrng string

	KeepLogDays int

	Parameters map[string]string
}
