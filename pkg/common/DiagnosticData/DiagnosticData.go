package mm

import "sync/atomic"

type DiagnosticData struct {
	// Number of all incoming requests.
	totalRequestsCount atomic.Uint64

	// Number of successfully finished requests.
	successfulRequestsCount atomic.Uint64
}

func (dd *DiagnosticData) GetTotalRequestsCount() (trc uint64) {
	return dd.totalRequestsCount.Load()
}

func (dd *DiagnosticData) IncTotalRequestsCount() {
	dd.totalRequestsCount.Add(1)
}

func (dd *DiagnosticData) GetSuccessfulRequestsCount() (src uint64) {
	return dd.successfulRequestsCount.Load()
}

func (dd *DiagnosticData) IncSuccessfulRequestsCount() {
	dd.successfulRequestsCount.Add(1)
}
