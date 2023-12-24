package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/MarlikAlmighty/mdns/internal/config"
	"github.com/MarlikAlmighty/mdns/internal/data"
)

func TestNew(t *testing.T) {
	// Test case 1: New core with valid resolved data and configuration
	resolvedData := &data.ResolvedData{}
	configuration := &config.Configuration{}
	core := New(resolvedData, configuration)
	assert.NotNil(t, core)
	assert.Equal(t, resolvedData, core.Resolver)
	assert.Equal(t, configuration, core.Config)

	// Test case 2: New core with nil resolved data and configuration
	core = New(nil, nil)
	assert.NotNil(t, core)
	assert.Nil(t, core.Resolver)
	assert.Nil(t, core.Config)
}
