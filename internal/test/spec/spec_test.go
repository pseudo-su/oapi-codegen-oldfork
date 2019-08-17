package spec

import (
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestTemplateInline(t *testing.T) {

	swaggerUI, err := GetSwaggerUI("/swagger/spec.json")
	assert.Equal(t, err, nil)
	err = cupaloy.SnapshotMulti("simple JSON", swaggerUI)
}
