package workload

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/model"
	"github.com/synadia-io/connect/runtime"
)

func runOutlet(ctx context.Context, runtime *runtime.Runtime, steps model.Steps) error {
	// we only support core consumers and nats sinks
	if steps.Consumer.JetStream != nil {
		return fmt.Errorf("FATAL: config: only core consumers are supported")
	}
	if steps.Sink.Type != "nats" {
		return fmt.Errorf("FATAL: sink type %s is not supported", steps.Sink.Type)
	}

	// connect to the sink nats
	sinkNc, err := connectToNats(steps.Sink.Config)
	if err != nil {
		return fmt.Errorf("FATAL: unable to connect to sink nats: %v", err)
	}
	defer sinkNc.Close()

	// connect to the nats environment containing the data we want to consumer
	dataNc, err := nats.Connect(steps.Consumer.NatsConfig.Url, steps.Consumer.NatsConfig.Opts()...)
	if err != nil {
		return fmt.Errorf("FATAL: unable to connect to the nats server at the data side: %v", err)
	}
	defer dataNc.Close()

	sinkSubject := ReadString(steps.Sink.Config, "subject")

	// the process is rather straightforward; we subscribe to the nats server containing the data and for each message
	// we receive, we publish it to the sink nats server
	_, err = dataNc.Subscribe(steps.Consumer.Subject, func(m *nats.Msg) {
		res := nats.NewMsg(sinkSubject)
		res.Data = m.Data
		res.Header = m.Header
		if err := sinkNc.PublishMsg(res); err != nil {
			runtime.Logger.Error(fmt.Sprintf("failed to publish message: %v", err))
		}
	})
	if err != nil {
		return fmt.Errorf("FATAL: unable to subscribe to the nats server at the data side: %v", err)
	}

	<-ctx.Done()
	return nil
}
