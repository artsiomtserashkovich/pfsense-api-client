package pfsenseapi

import (
	"errors"
	"time"
)

type Options func(o options)

func defaultOptions() options {
	return options{
		skipTLS: true,
		timeout: 5 * time.Second,
	}
}

type options struct {
	host     string
	username string
	password string
	skipTLS  bool
	timeout  time.Duration
}

func WithHost(host string) Options {
	return func(o options) {
		o.host = host
	}
}

func WithCredentials(username string, password string) Options {
	return func(o options) {
		o.username = username
		o.password = password
	}
}

func WithSkipTLS() Options {
	return func(o options) {
		o.skipTLS = true
	}
}

func WithTimeout(timeout time.Duration) Options {
	return func(o options) {
		o.timeout = timeout
	}
}

func (o options) validate() error {
	if o.host == "" {
		return errors.New("missing host")
	}
	if o.username == "" {
		return errors.New("missing username")
	}
	if o.password == "" {
		return errors.New("missing password")
	}
	return nil
}
