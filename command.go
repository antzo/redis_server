package redis_server

const separator = "\r\n"

const (
	SimpleStringType = "SimpleString"
	BulkStringType   = "BulkStrings"
	ArraysType       = "Arrays"
)

type Command struct {
	Type string
	Data string
}
