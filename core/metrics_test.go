package core_test

import (
	"mt-hosting-manager/core"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeExporterMetrics(t *testing.T) {
	data, err := os.ReadFile("testdata/metrics.txt")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	metrics, err := core.ParseNodeExporterMetrics(data)
	assert.NoError(t, err)
	assert.NotNil(t, metrics)

	assert.Equal(t, 0, metrics.LoadPercent)
	assert.Equal(t, int64(19809288192), metrics.DiskSize)
	assert.Equal(t, int64(17007599616), metrics.DiskUsed)
	assert.Equal(t, int64(2009100288), metrics.MemorySize)
	assert.Equal(t, int64(340316160), metrics.MemoryUsed)
}
