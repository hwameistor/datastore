package utils

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	log "github.com/sirupsen/logrus"
)

// BuildInClusterClientset builds a kubernetes in-cluster clientset
func BuildInClusterClientset() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.WithError(err).Fatal("Failed to build kubernetes config")
	}
	return kubernetes.NewForConfigOrDie(config)
}
