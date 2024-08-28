package vmprocessio

import "context"

type CoreVmClientProcess struct {
	cancel context.CancelFunc
}
