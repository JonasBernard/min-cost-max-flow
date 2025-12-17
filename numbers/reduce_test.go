package numbers_test

import (
	"testing"

	"github.com/JonasBernard/min-cost-max-flow/numbers"
	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	p, q := numbers.Reduce(0.33)
	assert.Equal(t, 33, p)
	assert.Equal(t, 100, q)
}
