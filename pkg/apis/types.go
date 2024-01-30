package apis

import "time"

type DataServer struct {
	Type      string
	Connected bool
	Endpoint  string

	SubDirs map[string]*DataDirectory

	AllDirs map[string]*DataDirectory
	AllObjs map[string]*DataObject
}

type DataDirectory struct {
	Name   string
	Path   string
	Prefix string // s3 bucket only

	SubDirs map[string]*DataDirectory
	Objects map[string]*DataObject
}

type DataObject struct {
	Name  string
	Path  string
	Size  int64
	MTime time.Time
}
