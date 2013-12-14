package factory

import (
	"../exchanges/bitcoincharts"
	"../exchanges/btce"
	core "../core"
	"errors"
	"os"
)

func NewExchange(exchange string) (core.Exchange, error) {
	return NewExchangeWithConfig(exchange, map[string]interface{}{})
}

func NewExchangeWithConfig(exchange string, config map[string]interface{}) (core.Exchange, error) {
	switch exchange {
	case "btce":
		exchange, err := btce.NewExchange(merge(config, map[string]interface{}{
			"key":    getEnv("BTCE_KEY", true),
			"secret": getEnv("BTCE_SECRET", true),
		}))
		if err != nil {
			return nil, err
		} else {
			return exchange, nil
		}
	case "bitcoincharts":
		return bitcoincharts.NewExchange(config), nil
	default:
		return nil, errors.New("Unknown exchange " + exchange)
	}

	return nil, nil
}

// get a key from the os env variables, optionally panic
func getEnv(key string, allowEmpty bool) string {
	value := os.Getenv(key)
	if !allowEmpty && value == "" {
		panic("ENV variable " + key + " needs to be set")
	}
	return value
}

// merge two maps, return union
func merge(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	for key, val := range m2 {
		m1[key] = val
	}
	return m1
}
