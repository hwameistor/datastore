package controller

import (
	"github.com/hwameistor/datastore/pkg/controller/storagebackend"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, storagebackend.Add)
}
