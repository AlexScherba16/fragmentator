package file_runner

import (
	"context"
	"fragmentator/internal/config"
	"fragmentator/internal/utils"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

type IRunner interface {
	Run()
	Name() string
}

type fileRunner struct {
	ctx      context.Context
	wg       *sync.WaitGroup
	filePath string
	conf     config.FilesConfig
}

func NewFileRunner(ctx context.Context, wg *sync.WaitGroup, filePath string, config config.FilesConfig) (*fileRunner, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}

	return &fileRunner{
		ctx:      ctx,
		wg:       wg,
		filePath: filePath,
		conf:     config,
	}, nil
}

func (f *fileRunner) Run() {
	defer f.wg.Done()

	file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_TRUNC, 666)
	if err != nil {
		log.Fatal("Can't open file ", "file", f.filePath, "error", err.Error())
		return
	}
	defer file.Close()

	chunkSize := f.conf.ChunkSizeMb() * 1024 * 1024
	timeOut := time.Duration(f.conf.IoIntervalMs())
	maxFileSize := int64(f.conf.MaxFileSizeMb() * 1024 * 1024)
	data := utils.GenerateRandomDataChunk(chunkSize)
	isAppendData := true

	for {
		select {
		case <-f.ctx.Done():
			log.Info("Shutdown file runner", "file", f.filePath)
			return
		default:
			if isAppendData {
				_, err = file.Write(data)
				if err != nil {
					log.Fatalf("Failed to write file: %s", err)
				}

				fileInfo, err := file.Stat()
				if err != nil {
					log.Fatalf("Failed to get file stat: %s", err)
				}

				if fileInfo.Size() >= maxFileSize {
					log.Debug("Swith to \"Truncate mode\"", "size", fileInfo.Size())
					isAppendData = false
				}

			} else {
				fileInfo, err := file.Stat()
				if err != nil {
					log.Fatalf("Failed to get file stat: %s", err)
				}

				truncatedSize := fileInfo.Size() - int64(chunkSize)
				if truncatedSize <= 0 {
					truncatedSize = 0
					log.Debug("Swith to \"Append mode\"", "file", fileInfo.Size())
					isAppendData = true
				}

				err = file.Truncate(truncatedSize)
				if err != nil {
					log.Fatalf("Failed to truncate file stat: %s", err)
				}

				if truncatedSize == 0 {
					_, err = file.Seek(0, 0)
					if err != nil {
						log.Fatalf("Failed to seek file: %s", err)
					}
				}

			}
			time.Sleep(timeOut * time.Millisecond)
		}
	}
}

func (f *fileRunner) Name() string {
	return f.filePath
}
