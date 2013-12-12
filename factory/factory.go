package factory

import (
	"errors"
	"os"
	"../btce"
	core "../core"
)

func NewExchange(exchange string) (core.Exchange, error) {
	switch exchange {
		case "btce":
			exchange, err := btce.NewExchange(map[string]string{
				"key": getEnv("BTCE_KEY", true),
				"secret": getEnv("BTCE_SECRET", true),
			})
			if err != nil {
				return nil, err
			} else {
				return exchange, nil
			}
		default:
			return nil, errors.New("Unknown exchange "+exchange)
	}

	return nil, nil
}

func getEnv(key string, allowEmpty bool) string {
	value := os.Getenv(key)
	if !allowEmpty && value == "" {
		panic("ENV variable " +key + " needs to be set")
	}
	return value
}