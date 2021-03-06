package pipeline

import (
	"github.com/sirupsen/logrus"
	"github.com/sunet/tq/pkg/message"
)

var Log = logrus.New()

type Pipeline func(...*message.MessageChannel) *message.MessageChannel

func CallPipeline(p Pipeline, cs ...*message.MessageChannel) *message.MessageChannel {
	return p(cs...)
}

func Merge(cs ...*message.MessageChannel) *message.MessageChannel {
	return message.ProcessChannels(func(v message.Message) (message.Message, error) {
		return v, nil
	}, "merge", cs...)
}

func LogMessages(cs ...*message.MessageChannel) *message.MessageChannel {
	return message.ProcessChannels(func(v message.Message) (message.Message, error) {
		m, err := message.FromJson(v)
		if err != nil {
			Log.Errorf("Unable to serialize json: %s", err.Error())
			return nil, err
		} else {
			Log.Print(string(m))
			return v, nil
		}
	}, "log", cs...)
}

func Run(cs ...*message.MessageChannel) {
	if len(cs) == 0 || cs == nil {
		cs = message.AllFinalChannels()
	}
	Log.Debugf("running %d final channels", len(cs))
	if len(cs) == 0 {
		select {} // for some reason the user wants us to block forever...
	} else if len(cs) == 1 {
		cs[0].Sink()
	} else {
		Merge(cs...).Sink()
	}
}
