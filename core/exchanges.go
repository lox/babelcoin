package babelcoin

import (
	"errors"
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
