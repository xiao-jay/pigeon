package plugins

type Plugin interface {
	SendMessage(msg interface{}) error
}
