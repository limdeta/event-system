package domain

type EventRegistry interface {
	ResolveChannel(channel string) (endpoint string, schemaName string, err error)
}
