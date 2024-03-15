package datamanager

type Manager interface {
	Cook() error

	Run(stopCh <-chan struct{})
}
