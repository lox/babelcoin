package babelcoin

import (
	"errors"
	"os"
	"strings"
)

var exchanges = make(map[string]ExchangeFactory)

// return an instance of an exchange, given it's string name
func NewExchange(key string, config map[string]interface{}) (Exchange, error) {
	parts := strings.SplitN(key, ":", 2)
	factory, ok := exchanges[parts[0]]
	if !ok {
		return nil, errors.New("No driver registered for " + parts[0])
	}

	return factory(key, config), nil
}

// called by drivers when initializing
func AddExchangeFactory(key string, factory ExchangeFactory) {
	exchanges[key] = factory
}

// read an exchanges key/secrets from env
func EnvExchangeConfig(key string) map[string]interface{} {
	config := map[string]interface{}{}
	for _, v := range os.Environ() {
		if strings.Index(v, strings.ToUpper(key)+"_") == 0 {
			parts := strings.SplitN(v, "=", 2)
			keyParts := strings.SplitN(parts[0], "_", 2)
			config[strings.ToLower(keyParts[1])] = parts[1]
		}
	}
	return config
}
