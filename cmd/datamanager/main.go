package main

import (
	"flag"
	"sync"

	"github.com/hwameistor/datastore/pkg/datamanager"

	log "github.com/sirupsen/logrus"
)

const (
	baseModelLocalPathDirDefault    = "/data/models"
	checkpointLocalPathDirDefault   = "/data/checkpoints"
	trainingdataLocalPathDirDefault = "/data/training"
	baseModelConfigNameDefault      = "test"
	checkpointConfigNameDefault     = "test"
	trainingdataConfigNameDefault   = "test"
)

var (
	nodeName  = flag.String("nodename", "localhost", "Node name")
	namespace = flag.String("namespace", "default", "Namespace of the Pod")
	isMaster  = flag.Bool("isTrainMaster", false, "role of the training Pod")

	// baseModelRemoteURL    = flag.String("baseModelRemoteURL", "", "Download URL of the base model, e.g. s3://localhost:9001/bucket/model.pth")
	// baseModelFileName     = flag.String("baseModelName", "", "File name of the base model, e.g. model.pth")
	// baseModeS3AccessKey   = flag.String("baseModelS3AccessKey", "", "AccessKey to access S3's bucket")
	// baseModeS3SecretKey   = flag.String("baseModelS3SecretKey", "", "SecretKey to access S3's bucket")
	baseModelLocalPathDir        = flag.String("baseModelLocalDir", baseModelLocalPathDirDefault, "Local dir of the base model, e.g. /data/models")
	checkpointLocalPathDir       = flag.String("checkpointLocalDir", checkpointLocalPathDirDefault, "Local dir of the checkpoints,  e.g. /data/checkpoints")
	trainingdataLocalPathDir     = flag.String("trainingdataLocalDir", trainingdataLocalPathDirDefault, "Local dir of the training data,  e.g. /data/training")
	checkpointHostLocalPathDir   = flag.String("checkpointLocalDirOnHost", "", "Local dir of the checkpoints on host")
	trainingdataHostLocalPathDir = flag.String("trainingdataLocalDirOnHost", "", "Local dir of the training data on host")

	baseModelConfigName    = flag.String("baseModelConfigName", baseModelConfigNameDefault, "Name of the base mode config")
	checkpointConfigName   = flag.String("checkpointConfigName", checkpointConfigNameDefault, "Name of the checkpoint config")
	trainingdataConfigName = flag.String("trainingdataConfigName", trainingdataConfigNameDefault, "Name of the training data config")

	isRoleforInit = flag.Bool("isInitRole", true, "role of the container")
)

func setupLogging() {
	log.SetLevel(log.DebugLevel)
}

func compileParametersForBaseModel() *datamanager.BaseModelParameters {
	return &datamanager.BaseModelParameters{
		Namespace:    *namespace,
		ConfigName:   *baseModelConfigName,
		LocalPathDir: *baseModelLocalPathDir,
	}
}

func compileParametersForCheckpoint() *datamanager.CheckpointParameters {
	return &datamanager.CheckpointParameters{
		Namespace:          *namespace,
		NodeName:           *nodeName,
		ConfigName:         *checkpointConfigName,
		LocalPathDir:       *checkpointLocalPathDir,
		LocalPathDirOnHost: *checkpointHostLocalPathDir,
	}
}

func compileParametersForTrainingData() *datamanager.TrainingDataParameters {
	return &datamanager.TrainingDataParameters{
		Namespace:          *namespace,
		NodeName:           *nodeName,
		ConfigName:         *trainingdataConfigName,
		LocalPathDir:       *trainingdataLocalPathDir,
		LocalPathDirOnHost: *trainingdataHostLocalPathDir,
	}
}

func main() {
	flag.Parse()
	setupLogging()

	if *isRoleforInit {
		runAsInitRole()
	} else {
		runAsManagerRole()
	}
}

func runAsInitRole() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		if err := datamanager.NewBaseModelManager(compileParametersForBaseModel()).Cook(); err != nil {
			log.WithError(err).Fatal("Base model is not ready")
		} else {
			log.Debug("Base model is ready")
		}
		wg.Done()
	}()

	if *isMaster {
		wg.Add(1)
		go func() {
			log.Debug("I am the master")
			if err := datamanager.NewCheckpointManager(compileParametersForCheckpoint()).Cook(); err != nil {
				log.WithError(err).Fatal("Checkpoint is not ready")
			} else {
				log.Debug("Checkpoint is ready")
			}
			wg.Done()
		}()
	} else {
		wg.Add(1)
		go func() {
			log.Debug("I am a worker")
			if err := datamanager.NewTrainingDataManager(compileParametersForTrainingData()).Cook(); err != nil {
				log.WithError(err).Fatal("Training data is not ready")
			} else {
				log.Debug("Training data is ready")
			}
			wg.Done()
		}()
	}

	wg.Wait()

}

func runAsManagerRole() {
	stopCh := make(chan struct{})

	go datamanager.NewBaseModelManager(compileParametersForBaseModel()).Run(stopCh)
	if *isMaster {
		go datamanager.NewCheckpointManager(compileParametersForCheckpoint()).Run(stopCh)
	} else {
		go datamanager.NewTrainingDataManager(compileParametersForTrainingData()).Run(stopCh)
	}

	<-stopCh
}
