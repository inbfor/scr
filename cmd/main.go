package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"ab_script/internal/config"
	"ab_script/internal/models"
	"ab_script/internal/utils"
)

func main() {

	configFilename := *flag.String("configLoc", "/etc/user-storage-basket/config.json", "конфигурация со списком vol'ов и дисков, на которых они находятся:")
	splitFilesCount := *flag.Int("splitFilesCount", 256, "количество частей, на которые необходимо разбивать архивы vol'ов")
	outputDir := *flag.String("outputDir", "/srv/opt-filelists", "директория, в которую необходимо класть листинги:")
	ioutil := *flag.Float64("ioutil", 40.0, "необходимая максимальная утилизация диска скриптом в %")
	batchPeriod := time.Duration(*flag.Int("batchPeriod", 20, "минимальная длительность непрерывного процесса операций листинга в мс")) * time.Millisecond
	flag.Parse()

	log.Println("Script started")

	cfg, err := config.ParseConfig(configFilename)
	if err != nil {
		log.Fatalf("Error parsing config file: %s", err)
	}

	var wg sync.WaitGroup
	for vols, disk := range cfg.RootMap {
		wg.Add(1)

		go func(disk string, vols string) {
			defer wg.Done()

			params := models.WalkDiskParams{
				Disk:            disk,
				Vol:             vols,
				Outdir:          outputDir,
				SplitFilesCount: splitFilesCount,
				Ioutil:          ioutil,
				BatchPeriod:     batchPeriod,
			}
			if err := utils.WalkDisk(params); err != nil {
				log.Printf("Error walking disk: %v", err)
			}
		}(disk, vols)
	}

	wg.Wait()
	log.Println("Script ended")
}
