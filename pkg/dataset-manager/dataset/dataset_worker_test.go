package dataset

import (
	"fmt"
	storageapis "k8s.io/api/storage/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/storage/v1/util"
	"reflect"
	"strings"
	"testing"
)

var (
	hddSC          = storageapis.StorageClass{ObjectMeta: v1.ObjectMeta{Name: "hwameistor-local-storage-hdd"}, Provisioner: "lvm.hwameistor.io", Parameters: map[string]string{"poolClass": "HDD"}}
	ssdSC          = storageapis.StorageClass{ObjectMeta: v1.ObjectMeta{Name: "hwameistor-local-storage-ssd"}, Provisioner: "lvm.hwameistor.io", Parameters: map[string]string{"poolClass": "SSD"}}
	nvmeSC         = storageapis.StorageClass{ObjectMeta: v1.ObjectMeta{Name: "hwameistor-local-storage-nvme"}, Provisioner: "lvm.hwameistor.io", Parameters: map[string]string{"poolClass": "NVMe"}}
	invalidOtherSC = storageapis.StorageClass{ObjectMeta: v1.ObjectMeta{Name: "other-sc"}, Provisioner: "storage.io"}
)

func TestSortHwameiStorageClasses(t *testing.T) {
	defaultHDDSC := hddSC
	defaultHDDSC.Annotations = map[string]string{"storageclass.kubernetes.io/is-default-class": "true"}

	testCases := []struct {
		desc     string
		input    []storageapis.StorageClass
		expected []storageapis.StorageClass
	}{
		{
			desc:     "one default storageclass(HDD) and three normal storageclasses",
			input:    []storageapis.StorageClass{defaultHDDSC, ssdSC, nvmeSC},
			expected: []storageapis.StorageClass{defaultHDDSC, nvmeSC, ssdSC},
		}, {
			desc:     "no default storageclass but three normal storageclasses",
			input:    []storageapis.StorageClass{hddSC, ssdSC, nvmeSC},
			expected: []storageapis.StorageClass{nvmeSC, ssdSC, hddSC},
		},
		{
			desc:     "one invalid storageclass and three normal storageclasses",
			input:    []storageapis.StorageClass{invalidOtherSC, hddSC, ssdSC, nvmeSC},
			expected: []storageapis.StorageClass{nvmeSC, ssdSC, hddSC},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.desc, func(t *testing.T) {
			sortedSC := sortHwameiStorageClasses(testCase.input)
			if !reflect.DeepEqual(testCase.expected, sortedSC) {
				t.Errorf("Expected: %v, got: %v", testCase.expected, sortedSC)
			}
			t.Logf("testcase output: %v", storageClassSlice(sortedSC))
			return
		})
	}
}

type storageClassSlice []storageapis.StorageClass

func (ss storageClassSlice) String() (so string) {
	for _, s := range ss {
		so += s.Name + "(" + fmt.Sprintf("%v", util.IsDefaultAnnotation(s.ObjectMeta)) + "," + s.Parameters["poolClass"] + ")" + ","
	}
	return strings.TrimSuffix(so, ",")
}
