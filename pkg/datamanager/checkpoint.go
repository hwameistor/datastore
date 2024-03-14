package datamanager

import (
	"context"
	"fmt"
	"os"
	"time"

	datastoreclientsetv1alpha1 "github.com/hwameistor/datastore/pkg/apis/client/clientset/versioned/typed/datastore/v1alpha1"
	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
	"github.com/hwameistor/datastore/pkg/storage/minio"
	"github.com/hwameistor/datastore/pkg/storage/ssh"
	"github.com/hwameistor/datastore/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	fswatcher "github.com/radovskyb/watcher"
	log "github.com/sirupsen/logrus"
)

type CheckpointParameters struct {
	NodeName           string
	Namespace          string
	ConfigName         string
	LocalPathDir       string
	LocalPathDirOnHost string
}

type CheckpointManager struct {
	params *CheckpointParameters

	config *datastorev1alpha1.CheckpointSpec

	clientset *datastoreclientsetv1alpha1.DatastoreV1alpha1Client
}

func NewCheckpointManager(params *CheckpointParameters) Manager {
	if err := utils.Mkdir(params.LocalPathDir); err != nil {
		log.WithField("localDir", params.LocalPathDir).WithError(err).Fatal("Failed to make local directory for checkpoints")
	}

	clientset := utils.BuildInClusterDataStoreClientset()
	cfg, err := clientset.Checkpoints(params.Namespace).Get(context.Background(), params.ConfigName, metav1.GetOptions{})
	if err != nil {
		log.WithField("config", params.ConfigName).WithError(err).Fatal("Failed to fetch config for checkpoints")
	}

	return &CheckpointManager{
		params:    params,
		clientset: clientset,
		config:    &cfg.Spec,
	}
}

func (mgr *CheckpointManager) getLatestRecord() *datastorev1alpha1.CheckpointRecord {

	cfg, err := mgr.clientset.Checkpoints(mgr.params.Namespace).Get(context.Background(), mgr.params.ConfigName, metav1.GetOptions{})
	if err != nil {
		log.WithField("config", mgr.params.ConfigName).WithError(err).Fatal("Failed to get config for checkpoints")
	}

	records := cfg.Status.Records

	var latestRecord *datastorev1alpha1.CheckpointRecord = nil
	for i := range records {
		if records[i].ExpiredTime != nil {
			continue
		}
		if latestRecord == nil {
			latestRecord = records[i]
			continue
		}
		if records[i].CreateTime.After(latestRecord.CreateTime.Time) {
			latestRecord = records[i]
		}
	}

	return latestRecord
}

func (mgr *CheckpointManager) isReady() bool {
	record := mgr.getLatestRecord()
	if record == nil {
		return true
	}

	files, err := os.ReadDir(mgr.params.LocalPathDir)
	if err != nil {
		log.WithField("dir", mgr.params.LocalPathDir).WithError(err).Fatal("Failed to access local directory for checkpoints")
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if f.Name() == record.Name {
			return true
		}
	}
	return false
}

func (mgr *CheckpointManager) Cook() error {
	if mgr.isReady() {
		return nil
	}

	ckpt := mgr.getLatestRecord()

	err := mgr.restoreFromPeer(ckpt)
	if err != nil {
		log.WithField("checkpoint", ckpt).WithError(err).Error("Failed to restore checkpoint from the peer")
		if ckpt.BackedUp {
			return mgr.restoreFromBackup(ckpt)
		}
	}
	return err
}

func (mgr *CheckpointManager) restoreFromPeer(ckpt *datastorev1alpha1.CheckpointRecord) error {
	log.WithField("checkpoint", ckpt).Debug("Restoring the checkpoint from a peer ...")
	fileLocalPath := fmt.Sprintf("%s/%s", mgr.params.LocalPathDir, ckpt.Name)
	fileRemotePath := fmt.Sprintf("%s/%s", ckpt.DirOnHost, ckpt.Name)

	return ssh.DownloadObject(ckpt.NodeName, fileRemotePath, fileLocalPath)
}

func (mgr *CheckpointManager) restoreFromBackup(ckpt *datastorev1alpha1.CheckpointRecord) error {
	log.WithField("checkpoint", ckpt).Debug("Restoring the checkpoint from the backup ...")

	if mgr.config.Backup == nil {
		log.Error("No backup configuration")
		return fmt.Errorf("no backup config")
	}
	if !ckpt.BackedUp {
		log.WithField("checkpoint", ckpt.Name).Debug("no backup found")
		return fmt.Errorf("no backup found")
	}
	if mgr.config.Backup.Proto == "minio" && mgr.config.Backup.MinIO != nil {
		return minio.DownloadObject(mgr.config.Backup.MinIO, ckpt.Name, fmt.Sprintf("%s/%s", mgr.params.LocalPathDir, ckpt.Name))
	}
	log.WithField("config", mgr.config.Backup).Error("Invalid backup configration")
	return fmt.Errorf("invalid backup config")
}

func (mgr *CheckpointManager) Run(stopCh <-chan struct{}) {
	logCtx := log.WithFields(log.Fields{"directory": mgr.params.LocalPathDir})
	logCtx.Debug("Watching for the new checkpoint ...")
	// creates a new file watcher
	watcher := fswatcher.New()
	watcher.FilterOps(fswatcher.Create)

	if err := watcher.Add(mgr.params.LocalPathDir); err != nil {
		logCtx.WithError(err).Fatal("Failed to watch checkpoints directory")
	}

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Event:
				log.WithFields(log.Fields{"checkpoint": event}).Debug("Captured a new checkpoint")
				go mgr.handleNewCheckpoint(&event)

				// watch for errors
			case err := <-watcher.Error:
				log.WithFields(log.Fields{"error": err.Error()}).Fatal("Got a watch error")

			case <-watcher.Closed:
				return
			case <-stopCh:
				watcher.Close()
			}
		}
	}()

	if err := watcher.Start(1 * time.Second); err != nil {
		logCtx.WithError(err).Fatal("Failed to start the watcher")
	}

	<-stopCh
}

func (mgr *CheckpointManager) handleNewCheckpoint(event *fswatcher.Event) {
	logCtx := log.WithField("checkpoint", event.Path)
	logCtx.Debug("Start to handle the new checkpint ...")

	var lastsize int64 = -1
	var currsize int64 = 0
	for currsize > lastsize {
		lastsize = currsize
		time.Sleep(1 * time.Second)

		f, err := os.Stat(event.Path)
		if err != nil {
			logCtx.WithError(err).Fatal("Failed to access the checkpoint")
		}
		currsize = f.Size()
		logCtx.WithFields(log.Fields{"curr": currsize, "last": lastsize}).Debug("Got the size of checkpoint")
	}

	checksum, err := utils.Checksum(event.Path)
	if err != nil {
		logCtx.WithError(err).Error("Failed to calculate checksum")
	}
	rec := &datastorev1alpha1.CheckpointRecord{
		Name:       event.Name(),
		NodeName:   mgr.params.NodeName,
		DirOnHost:  mgr.params.LocalPathDirOnHost,
		Size:       fmt.Sprintf("%d", currsize),
		Checksum:   checksum,
		CreateTime: &metav1.Time{Time: event.ModTime()},
	}

	if mgr.config.Backup != nil {
		if err := mgr.backupCheckpoint(event.Path, event.Name()); err != nil {
			logCtx.WithError(err).Error("Failed to backup the checkpoint")
			rec.BackedUp = false
		} else {
			rec.BackedUp = true
			logCtx.Debug("Successfully backup the checkpoint")
		}
	}

	if err := mgr.recordCheckpoint(rec); err != nil {
		logCtx.WithError(err).Error("Failed to record the checkpoint")
	}
}

func (mgr *CheckpointManager) backupCheckpoint(fpath string, fname string) error {
	log.WithField("checkpoint", fpath).Debug("Start to backup the checkpoint ...")

	if mgr.config.Backup.Proto == "minio" && mgr.config.Backup.MinIO != nil {
		return minio.UploadObject(mgr.config.Backup.MinIO, fname, fpath)
	}
	log.WithField("config", mgr.config.Backup).Error("Invalid backup configration")
	return fmt.Errorf("invalid backup config")
}

func (mgr *CheckpointManager) recordCheckpoint(rec *datastorev1alpha1.CheckpointRecord) error {
	log.WithField("record", rec).Debug("Saving the checkpoint record ...")

	ctx := context.Background()
	ckpt, err := mgr.clientset.Checkpoints(mgr.params.Namespace).Get(ctx, mgr.params.ConfigName, metav1.GetOptions{})
	if err != nil {
		log.WithField("config", mgr.params.ConfigName).WithError(err).Fatal("Failed to get config for checkpoints")
	}

	ckpt.Status.Records = append([]*datastorev1alpha1.CheckpointRecord{rec}, ckpt.Status.Records...)
	_, err = mgr.clientset.Checkpoints(mgr.params.Namespace).UpdateStatus(ctx, ckpt, metav1.UpdateOptions{})
	return err
}
