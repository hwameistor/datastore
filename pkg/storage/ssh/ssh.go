package ssh

import (
	"github.com/pytool/ssh"
)

func DownloadObject(host string, remoteFilePath string, localFilePath string) error {

	conf := ssh.Config{
		Host: host,
		Port: 22,
	}
	client, err := ssh.New(&conf)
	if err != nil {
		return err
	}
	defer client.Close()

	return client.Download(remoteFilePath, localFilePath)
}
