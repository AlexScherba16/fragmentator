package config

import (
	"encoding/json"
	"os"
)

type FilesConfig struct {
	files       uint32
	ioInterval  uint32
	chunkSize   uint32
	maxFileSize uint32
}

func (c *FilesConfig) FilesNumber() uint32   { return c.files }
func (c *FilesConfig) IoIntervalMs() uint32  { return c.ioInterval }
func (c *FilesConfig) ChunkSizeMb() uint32   { return c.chunkSize }
func (c *FilesConfig) MaxFileSizeMb() uint32 { return c.maxFileSize }

type AppConfig struct {
	filesConfig   *FilesConfig
	directories   []string
	launchTimeout uint32
	logLevel      string
}

func (c *AppConfig) FilesContext() *FilesConfig { return c.filesConfig }
func (c *AppConfig) Dirs() []string             { return c.directories }
func (c *AppConfig) LaunchTimeoutMs() uint32    { return c.launchTimeout }
func (c *AppConfig) LogLevel() string           { return c.logLevel }

func NewAppConfig(configPath string) (*AppConfig, error) {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Anonymous structs
	type fileJsonConfig struct {
		Files       uint32 `json:"files"`
		IoInterval  uint32 `json:"io_timeout_ms"`
		ChunkSize   uint32 `json:"chunk_size_mb"`
		MaxFileSize uint32 `json:"max_file_size_mb"`
	}
	type appJsonConfig struct {
		FileContext   fileJsonConfig `json:"files_config"`
		Directories   []string       `json:"directories"`
		LaunchTimeout uint32         `json:"thread_launch_timeout_ms"`
		LogLevel      string         `json:"log_level"`
	}

	conf := &appJsonConfig{}
	err = json.Unmarshal(file, conf)
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		filesConfig: &FilesConfig{
			files:       conf.FileContext.Files,
			ioInterval:  conf.FileContext.IoInterval,
			chunkSize:   conf.FileContext.ChunkSize,
			maxFileSize: conf.FileContext.MaxFileSize,
		},
		directories:   conf.Directories,
		launchTimeout: conf.LaunchTimeout,
		logLevel:      conf.LogLevel,
	}, nil
}
