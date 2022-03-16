// Package command contains cli commands.
package command

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	gzGlue "github.com/gozix/glue/v2"
	gzViper "github.com/gozix/viper/v2"
	"github.com/sarulabs/di/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrKeyNotExist = fmt.Errorf("config path not exist")
)

type (
	downloader struct {
		logger *log.Logger
		config *viper.Viper
	}
)

// DefCommandDownloadName is container name.
const DefCommandDownloadName = "cli.command.download"

// DefCommandDownload register command in di container.
func DefCommandDownload() di.Def {
	return di.Def{
		Name: DefCommandDownloadName,
		Tags: []di.Tag{{
			Name: gzGlue.TagCliCommand,
		}},
		Build: func(ctn di.Container) (_ interface{}, err error) {
			return &cobra.Command{
				Use:           "download",
				Short:         "download imdb data",
				SilenceUsage:  true,
				SilenceErrors: true,
				RunE:          DownloadRunE(ctn),
			}, nil
		},
	}
}

func DownloadRunE(ctn di.Container) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		var (
			config = ctn.Get(gzViper.BundleName).(*viper.Viper)
			logger = log.New(os.Stdout, "", log.LstdFlags)
		)
		config = config.Sub("app.imdb")
		if config == nil {
			logger.Println("app.imdb not exist")
			return ErrKeyNotExist
		}

		var d = downloader{
			logger: logger,
			config: config,
		}

		return d.Handler(cmd, args)
	}
}

// Handler run.
func (s *downloader) Handler(_ *cobra.Command, _ []string) (err error) {

	var host = s.config.GetString("http_addr")
	if host == "" {
		s.logger.Println("http_addr not exist")
		return ErrKeyNotExist
	}

	var files = s.config.GetStringMapString("files")
	if files == nil {
		s.logger.Println("host not exist")
		return ErrKeyNotExist
	}

	var downloadPath = s.config.GetString("download_path")
	if files == nil {
		s.logger.Println("download_path not exist")
		return ErrKeyNotExist
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		s.logger.Printf("pwd: %s", err)
	}

	s.logger.Println("Start...")
	var wg sync.WaitGroup
	for name, filename := range files {
		wg.Add(1)
		go func(name, host, filename string) {
			defer wg.Done()
			if err := s.downloadFile(
				fmt.Sprintf("%s/%s/%s", dir, downloadPath, filename),
				host+filename,
			); err != nil {
				s.logger.Printf("download %s failed: %s", filename, err)
			}
			s.logger.Printf("%s was downloaded", filename)
			if err := s.escapeQuote(
				fmt.Sprintf("%s/%s/%s", dir, downloadPath, filename),
			); err != nil {
				s.logger.Printf("quote escaped %s failed: %s", filename, err)
			}
			s.logger.Printf("%s was quote escaped", filename)
		}(name, host, filename)
	}
	wg.Wait()
	s.logger.Println("Finished!")

	return
}

func (s *downloader) downloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func (s *downloader) escapeQuote(filepath string) (err error) {
	cmd := fmt.Sprintf(
		`zcat < %[1]s | sed 's/\//\\\\\//g;s/"/\/"/g' | gzip -c > %[1]s.tmp && mv %[1]s.tmp %[1]s`,
		filepath)
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return err
	}
	return nil
}
