package main

import (
	"encoding/json"
	"flag"
	"fmt"
	filestore "lostpets/internal/data/file-store"
	"lostpets/internal/data/postgres"
	"lostpets/internal/http"
	"lostpets/internal/logging"
	"os"
)

type config struct {
	Logger    logging.LogrusConfig
	Server    http.Config      `json:"server"`
	Database  postgres.Config  `json:"db"`
	FileStore filestore.Config `json:"fileStore"`
}

var (
	version   string
	timestamp string
)

func main() {
	printVersion := flag.Bool("version", false, "print version and exit")
	configFileName := flag.String("c", "", "configuration file to use")
	flag.Parse()

	if *printVersion {
		fmt.Printf("Version: %s - %s", version, timestamp)
		os.Exit(0)
	}

	var config config
	err := config.load(*configFileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log, err := logging.NewLogrusWrapper(config.Logger)
	if err != nil {
		fmt.Printf("Failed to create logger: %s", err)
		os.Exit(1)
	}

	db, err := postgres.NewDBConnection(config.Database)
	if err != nil {
		fmt.Printf("Failed to create db: %s", err)
		os.Exit(1)
	}

	fs, err := filestore.NewFileStore(config.FileStore)
	if err != nil {
		fmt.Printf("Failed to create file store: %s", err)
		os.Exit(1)
	}

	http.StartServer(config.Server, db, db, fs, log, version)

}

func (c *config) load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(c)
	return err
}
