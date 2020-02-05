package income

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoneyScan(t *testing.T) {
	f := 1234.56
	m := new(Money)
	assert.NoError(t, m.Scan(f))
	assert.Equal(t, FromFloat(f), *m)
}
