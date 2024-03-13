package ftp

import (
	"io"
	"os"
	"time"

	"github.com/jlaffaye/ftp"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
)

func DownloadObject(spec *datastorev1alpha1.FTPSpec, fileName string, localFilePath string) error {

	c, err := ftp.Dial(spec.Endpoint, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return err
	}

	err = c.Login(spec.LoginUser, spec.LoginPassword)
	if err != nil {
		return err
	}

	c.ChangeDir(spec.Dir)

	res, err := c.Retr(fileName)
	if err != nil {
		return err
	}
	defer res.Close()

	outFile, err := os.Create(localFilePath)
	if err != nil {
		return err
	}

	defer outFile.Close()

	_, err = io.Copy(outFile, res)
	if err != nil {
		return err
	}

	return nil

}
