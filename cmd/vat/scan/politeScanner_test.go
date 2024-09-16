package scan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoliteScanner_Scan(t *testing.T) {
	politeScanner := &PoliteScanner{
		Host: "192.127.0.1",
	}
	res, err := politeScanner.Scan()

	assert.Nil(t, res)
	assert.NotNil(t, err)
}
