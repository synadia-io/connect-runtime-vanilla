package workload

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/model"
	"github.com/synadia-io/connect/runtime"
)

func runInlet(ctx context.Context, runtime *runtime.Runtime, steps model.Steps) error {
	// we only support core producers and nats sources
	if steps.Producer.JetStream != nil {
		return fmt.Errorf("FATAL: config: only core producers are supported")
	}
	if steps.Source.Type != "nats" {
		return fmt.Errorf("FATAL: source type %s is not supported", steps.Source.Type)
	}

	// connect to the source nats
	sourceNc, err := connectToNats(steps.Source.Config)
	if err != nil {
		return fmt.Errorf("FATAL: unable to connect to source nats: %v", err)
	}
	defer sourceNc.Close()

	// connect to the nats environment to which we want to send the data
	dataNc, err := nats.Connect(steps.Producer.NatsConfig.Url, steps.Producer.NatsConfig.Opts()...)
	if err != nil {
		return fmt.Errorf("FATAL: unable to connect to nats: %v", err)
	}
	defer dataNc.Close()

	sourceSubject := ReadString(steps.Source.Config, "subject")

	// the process is rather straightforward; we subscribe to the source nats and for each message we receive,
	// we publish it to the data nats server
	_, err = sourceNc.Subscribe(sourceSubject, func(m *nats.Msg) {
		res := nats.NewMsg(steps.Producer.Subject)
		res.Data = m.Data
		res.Header = m.Header

		if err := dataNc.PublishMsg(res); err != nil {
			runtime.Logger.Error(fmt.Sprintf("failed to publish message: %v", err))
			return
		}

		runtime.Logger.Debug(fmt.Sprintf("published message to %s", steps.Producer.Subject))
	})
	if err != nil {
		return fmt.Errorf("FATAL: unable to subscribe to source nats: %v", err)
	}

	<-ctx.Done()
	return nil
}
