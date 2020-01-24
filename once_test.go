package cancellation

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestOnce(t *testing.T) {
	once := NewOnce()
	calls := rand.Intn(100)
	doneCh := make(chan struct{}, calls)

	for i := 0; i < calls; i++ {
		go func(i int) {
			once.Cancel(errors.New(fmt.Sprintf("err%d", i)))
			doneCh <- struct{}{}
		}(i)
	}

	err := once.Wait()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "err")
	for i := 0; i < cap(doneCh); i++ {
		<-doneCh
	}
}
