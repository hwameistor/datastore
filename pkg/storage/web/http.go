package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	datastorev1alpha1 "github.com/hwameistor/datastore/pkg/apis/datastore/v1alpha1"
)

func DownloadObject(spec *datastorev1alpha1.HTTPSpec, fileName string, localFilePath string) error {

	addr := fmt.Sprintf("%s/%s", strings.TrimRight(spec.Url, "/"), fileName)

	_, err := os.Stat(localFilePath)
	if err == nil {
		// already exist
		return nil
	}

	out, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(addr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
