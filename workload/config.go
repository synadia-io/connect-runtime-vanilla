package workload

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func connectToNats(config map[string]any) (*nats.Conn, error) {
	// -- check if the url is set
	url := ReadString(config, "url")
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

	var opts []nats.Option
	authEnabled := readBool(config, "auth")
	if authEnabled {
		userJwt := ReadString(config, "user_jwt")
		userSeed := ReadString(config, "user_seed")

		if userJwt == "" || userSeed == "" {
			return nil, fmt.Errorf("user_jwt and user_seed are required in case of a secured nats server")
		}

		opts = append(opts, nats.UserJWTAndSeed(userJwt, userSeed))
	}

	return nats.Connect(url, opts...)
}

func readBool(config map[string]any, key string) bool {
	val, ok := config[key]
	if !ok {
		return false
	}

	return val.(bool)
}

func ReadString(config map[string]any, key string) string {
	val, ok := config[key]
	if !ok {
		return ""
	}

	return val.(string)
}
