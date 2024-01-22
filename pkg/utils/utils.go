package utils

import (
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	log "github.com/sirupsen/logrus"
)

// BuildInClusterClientset builds a kubernetes in-cluster clientset
func BuildInClusterClientset() *kubernetes.Clientset {
	config, err := config.GetConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to build kubernetes config")
	}
	return kubernetes.NewForConfigOrDie(config)
}
