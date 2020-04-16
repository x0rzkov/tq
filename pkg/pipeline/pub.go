package pipeline

import (
	"github.com/SUNET/tq/pkg/message"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol/pub"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

func MakePublishPipeline(url string) Pipeline {
	var err error
	var sock mangos.Socket
	var data []byte
	if sock, err = pub.NewSocket(); err != nil {
		Log.Panicf("can't create pub socket: %s", err.Error())
	}
	if err = sock.Listen(url); err != nil {
		Log.Panicf("can't listen to pub %s on socket: %s", url, err.Error())
	}

	return func(cs ...*message.MessageChannel) *message.MessageChannel {
		return message.ProcessChannels(func(o message.Message) (message.Message, error) {
			data, err = message.FromJson(o)
			if err != nil {
				Log.Errorf("Error serializing json: %s", err)
				return nil, err
			} else {
				sock.Send(data)
				return o, nil
			}
		}, cs...)
	}
}
