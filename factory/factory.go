package factory

import (
	"strings"
	"errors"
	"os"
	"../btce"
	core "../core"
)

func NewExchange(pair string) (core.Exchange, error) {
	parts := strings.Split(pair, "/")
	if len(parts) != 2 {
		return nil, errors.New("Expected pair in form of exchange/cur_cur")
	}

	switch parts[0] {
		case "btce":
			exchange, err := btce.NewExchange(parts[1], map[string]string{
				"key": getEnv("BTCE_KEY", true),
				"secret": getEnv("BTCE_SECRET", true),
			})
			if err != nil {
				return nil, err
			} else {
				return exchange, nil
			}
		default:
			return nil, errors.New("Unknown exchange "+parts[0])
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