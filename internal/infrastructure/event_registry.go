package infrastructure

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type EventChannelInfo struct {
	Endpoint   string `json:"endpoint"`
	SchemaName string `json:"schema"`
	Type       string `json:"type"`
}

type EventRegistry struct {
	channels map[string]EventChannelInfo
	filePath string
	mu       sync.RWMutex
}

func NewEventRegistryFromFile(path string) (*EventRegistry, error) {
	reg := &EventRegistry{
		filePath: path,
	}
	if err := reg.Reload(); err != nil {
		return nil, err
	}
	return reg, nil
}

// Reload перечитывает файл конфигурации каналов.
func (r *EventRegistry) Reload() error {
	f, err := os.Open(r.filePath)
	if err != nil {
		return fmt.Errorf("cannot open event registry config: %w", err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	chMap := make(map[string]EventChannelInfo)
	if err := dec.Decode(&chMap); err != nil {
		return fmt.Errorf("cannot decode event registry config: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.channels = chMap
	fmt.Println("[event-registry] config reloaded")
	return nil
}

func (r *EventRegistry) ResolveChannel(channel string) (string, string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	info, ok := r.channels[channel]
	if !ok {
		return "", "", fmt.Errorf("channel %q not found in event registry", channel)
	}
	return info.Endpoint, info.SchemaName, nil
}

// GetAllChannels возвращает копию текущей карты каналов (для просмотра через API).
func (r *EventRegistry) GetAllChannels() map[string]EventChannelInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make(map[string]EventChannelInfo, len(r.channels))
	for k, v := range r.channels {
		result[k] = v
	}
	return result
}
