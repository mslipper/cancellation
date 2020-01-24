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
	doneCh := make(chan struct{})
	go func() {
		for i := 0; i < rand.Intn(100); i++ {
			once.Cancel(errors.New(fmt.Sprintf("err%d", i)))
			doneCh <- struct{}{}
		}
	}()

	err := once.Wait()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "err0")
	<-doneCh
}
