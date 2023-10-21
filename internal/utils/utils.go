package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dsnet/compress/bzip2"

	"ab_script/internal/models"
)

const (
	offsetBasis uint32 = 0x811c9dc5
	prime       uint32 = 0x01000193
)

func WalkDisk(params models.WalkDiskParams) error {
	volDir := filepath.Join(params.Outdir, params.Vol)

	fileInfo, err := os.Stat(filepath.Join(params.Disk, params.Vol))
	if err != nil {
		return fmt.Errorf("can't read file's filesystem stats: %w", err)
	}

	timeFolderUnix := fileInfo.ModTime().UnixNano()
	timeFolderRFC := time.Unix(0, timeFolderUnix).Format(time.RFC3339)

	log.Println(volDir)
	if err := os.MkdirAll(volDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", volDir, err)
	}

	slcBzip2, slcFile := createFiles(filepath.Join(params.Outdir, params.Vol), params.Vol, params.SplitFilesCount)

	for i := range slcFile {
		defer slcFile[i].Close()
		defer slcBzip2[i].Close()
	}

	log.Printf("walk %s/%s: started", params.Disk, params.Vol)

	if err := filepath.WalkDir(filepath.Join(params.Disk, params.Vol), func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("prevent panic by handling failure accessing a path: %w", err)
		}

		infoFile, err := info.Info()
		if err != nil {
			return fmt.Errorf("some file error")
		}

		timeFileRFC := time.Unix(0, infoFile.ModTime().UnixNano()).Format(time.RFC3339)

		if info.IsDir() {
			return nil
		}

		index := hashByData([]byte(infoFile.Name())) % uint32(params.SplitFilesCount)
		bz := slcBzip2[index]

		row := []string{
			timeFolderRFC,
			timeFileRFC,
			strconv.FormatInt(infoFile.Size(), 10),
			strings.Join([]string{params.Vol, info.Name()}, "/"),
		}

		_, errBz := bz.Write([]byte(strings.Trim(strings.Join(row, ";"), ";") + "\n"))

		if errBz != nil {
			return errBz
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error walking filepath: %w", err)
	}

	log.Printf("walk %s/%s: ended", params.Disk, params.Vol)
	return nil
}

func hashByData(data []byte) uint32 {

	var hval = offsetBasis
	for _, b := range data {
		hval ^= uint32(b)
		hval *= prime
	}
	return hval
}

func createFiles(path string, vol string, numberOfFiles int) ([]*bzip2.Writer, []*os.File) {

	var slcBzip2 []*bzip2.Writer
	var sliceFile []*os.File
	var i int64

	for i = 0; i < int64(numberOfFiles); i++ {

		hex := strconv.FormatInt(i, 16)

		file, err := os.Create(filepath.Join(path, vol+"-"+hex+".bz2"))
		if err != nil {
			log.Println(err)
		}
		sliceFile = append(sliceFile, file)

		bz, err := bzip2.NewWriter(file, &bzip2.WriterConfig{Level: bzip2.BestCompression})

		if err != nil {
			log.Println(err)
		}

		slcBzip2 = append(slcBzip2, bz)
	}

	return slcBzip2, sliceFile
}
