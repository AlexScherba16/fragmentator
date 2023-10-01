package factory

import (
	"context"
	"fmt"
	"fragmentator/internal/config"
	"fragmentator/internal/runners/file_runner"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"sync"
)

func NewRunners(ctx context.Context, wg *sync.WaitGroup, dirPath string, config *config.AppConfig) ([]file_runner.IRunner, error) {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return nil, err
		}
		log.Debug("Created", "directory", dirPath)
	} else {
		log.Debug("Exists", "directory", dirPath)
	}

	runners := []file_runner.IRunner{}
	prefix := filepath.Base(dirPath)
	for i := uint32(0); i < config.FilesContext().FilesNumber(); i++ {
		fileName := fmt.Sprintf("%s_%d.txt", prefix, i)
		filePath := filepath.Join(dirPath, fileName)

		runner, err := file_runner.NewFileRunner(ctx, wg, filePath, *config.FilesContext())
		if err != nil {
			return nil, err
		}
		runners = append(runners, runner)
	}

	return runners, nil
}
