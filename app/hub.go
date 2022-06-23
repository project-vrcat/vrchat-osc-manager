package app

import "vrchat-osc-manager/internal/pubsub"

type (
	Message struct {
		method string
		avatar string
		parameter
	}
	parameter struct {
		name  string
		value []any
	}
)

var hub = pubsub.New[Message](10)
