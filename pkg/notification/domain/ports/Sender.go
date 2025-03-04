package ports

type Sender interface {
	Send(title string, message string) (interface{}, error)
}
