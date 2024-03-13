package datamanager

import (
	log "github.com/sirupsen/logrus"
)

type TrainingDataParameters struct {
	Namespace          string
	ConfigName         string
	NodeName           string
	LocalPathDir       string
	LocalPathDirOnHost string
}

type TrainingDataManager struct {
	params *TrainingDataParameters
}

func NewTrainingDataManager(params *TrainingDataParameters) Manager {
	return &TrainingDataManager{
		params: params,
	}
}

func (mgr *TrainingDataManager) isReady() bool {
	return true
}

func (mgr *TrainingDataManager) Cook() error {
	if mgr.isReady() {
		return nil
	}

	return nil
}

func (mgr *TrainingDataManager) Run(stopCh <-chan struct{}) {
	log.Debug("Watching for training data request ...")
}
