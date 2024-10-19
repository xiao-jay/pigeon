package channels

import "pigeon/config"

type Channel interface {
	SendMessage(msgs []config.Msg, extra any) error
}
