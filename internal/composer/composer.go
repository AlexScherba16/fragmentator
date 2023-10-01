package composer

import (
	"context"
	"fragmentator/internal/cli"
	"fragmentator/internal/config"
	"fragmentator/internal/constants"
	"fragmentator/internal/runners/factory"
	"fragmentator/internal/runners/file_runner"
	"github.com/charmbracelet/log"
	"sync"
	"time"
)

type composer struct {
	cancel          context.CancelFunc
	wg              *sync.WaitGroup
	runners         []file_runner.IRunner
	launcherTimeout uint32
}

func NewComposer() (*composer, error) {
	cli, err := cli.NewFlags()
	if err != nil {
		return nil, err
	}

	conf, err := config.NewAppConfig(cli.ConfigPath())
	if err != nil {
		return nil, err
	}
	if conf.LogLevel() == constants.MaxLogLevel {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	runners := []file_runner.IRunner{}

	for _, runDirectory := range conf.Dirs() {
		fileRunners, err := factory.NewRunners(ctx, wg, runDirectory, conf)
		if err != nil {
			return nil, err
		}
		log.Info("Create runners for", "directory", runDirectory)
		runners = append(runners, fileRunners...)
	}

	return &composer{
		cancel:          cancel,
		wg:              wg,
		runners:         runners,
		launcherTimeout: conf.LaunchTimeoutMs(),
	}, nil
}

func (c *composer) RunApplication() {
	timeOut := time.Duration(c.launcherTimeout)
	c.wg.Add(len(c.runners))

	for _, runner := range c.runners {
		log.Info("Launch", "runner", runner.Name())
		go runner.Run()
		time.Sleep(timeOut * time.Millisecond)
	}
}

func (c *composer) StopApplication() {
	c.cancel()
	c.wg.Wait()
}
