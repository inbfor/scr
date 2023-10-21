package models

import "time"

type DataRow struct {
	ParentMtime float64
	Mtime       float64
	Size        int64
	Disk        string
	ParentDir   string
	Name        string
}

type WalkDiskParams struct {
	Disk            string
	Vol             string
	Outdir          string
	SplitFilesCount int
	Ioutil          float64
	BatchPeriod     time.Duration
}
