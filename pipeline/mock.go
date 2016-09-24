package pipeline

import (
	"github.com/danielkrainas/csense/context"
)

type MockPipeline struct {
	SendFunc func(ctx context.Context, m Message)
}

var _ Pipeline = &MockPipeline{}

func (pipe *MockPipeline) Send(ctx context.Context, m Message) {
	if pipe.SendFunc != nil {
		pipe.Send(ctx, m)
	} else {
		logSender(ctx, m)
	}
}

func logSender(ctx context.Context, m Message) {
	context.GetLogger(ctx).Infof("sending message %+#v", m)
}
