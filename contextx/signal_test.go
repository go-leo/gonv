package contextx_test

import (
	"syscall"
	"testing"

	"github.com/go-leo/gox/contextx"
)

func TestSignal(t *testing.T) {
	ctx, cancel := contextx.Signal(syscall.SIGHUP)
	t.Log(ctx)
	t.Log(cancel)
}
