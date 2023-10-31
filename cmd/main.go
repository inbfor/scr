package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/scr/abscript/pkg/config"
	"github.com/scr/abscript/pkg/models"
	"github.com/scr/abscript/pkg/utils"
)

func main() {
	configFilename := flag.String("configLoc",
		"/etc/user-storage-basket/config.json",
		"конфигурация со списком vol'ов и дисков, на которых они находятся:")
	splitFilesCount := flag.Int(
		"splitFilesCount",
		256,
		"количество частей, на которые необходимо разбивать архивы vol'ов")
	outputDir := flag.String(
		"outputDir",
		"/srv/opt-filelists",
		"директория, в которую необходимо класть листинги:")
	ioutil := flag.Float64(
		"ioutil",
		40.0,
		"необходимая максимальная утилизация диска скриптом в %")

	batchPeriodML := flag.Int("batchPeriod",
		20,
		"минимальная длительность непрерывного процесса операций листинга в мс")
	flag.Parse()

	batchPeriod := time.Duration(
		*batchPeriodML) * time.Millisecond

	log.Println("Script started")

	cfg, err := config.ParseConfig(*configFilename)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	var waitGroup sync.WaitGroup
	for vols, disk := range cfg.RootMap {
		waitGroup.Add(1)

		go func(disk string, vols string) {
			defer waitGroup.Done()

			params := models.WalkDiskParams{
				Disk:            disk,
				Vol:             vols,
				Outdir:          *outputDir,
				SplitFilesCount: *splitFilesCount,
				Ioutil:          *ioutil,
				BatchPeriod:     batchPeriod,
			}
			if err := utils.WalkDisk(params); err != nil {
				log.Printf("Error walking disk: %v", err)
			}
		}(disk, vols)
	}

	waitGroup.Wait()
	log.Println("Script ended")
}
