package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq" //postgres driver import
	goose "github.com/pressly/goose"
)

const (
	postgresDBString = "user=%s password=%s host=%s port=%d dbname=%s sslmode=%s"
)

type (
	//this config overlaps with the db portion of the server config, so they can use the same config file if wanted

	dbConfig struct {
		MigrationPath string `json:"migrationPath"`
		Username      string `json:"migrationUser"`
		Password      string `json:"migrationPassword"`
		Host          string `json:"host"`
		Port          int    `json:"port"`
		DBName        string `json:"dbName"`
		SSLMode       string `json:"sslMode"`
	}

	config struct {
		DB dbConfig `json:"db"`
	}
)

func main() {
	// -- Flags -- //
	flags := flag.NewFlagSet("goose", flag.ExitOnError)
	configFile := flags.String("c", "./db.json", "config file path")
	printVersion := flags.Bool("version", false, "print version and exit")
	verbose := flags.Bool("v", false, "enable verbose mode") //just copying GOOSE Flags so can document them
	help := flags.Bool("h", false, "print help")
	flags.Usage = usage(flags)
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 1 || *help {
		flags.Usage()
		return
	}

	if *printVersion {
		fmt.Printf("Goose Version: %s\n", goose.VERSION)
		os.Exit(0)
	}

	if *verbose {
		goose.SetVerbose(true)
	}

	// load conf file
	var conf config
	err := conf.load(*configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading configuration file:", err)
		os.Exit(1)
	}

	command := args[0]
	dbString := conf.DB.postgresDBString()

	db, err := goose.OpenDBWithDriver("postgres", dbString)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	//migrations path is relative to the config
	migrationPath := conf.DB.MigrationPath
	if !filepath.IsAbs(migrationPath) {
		migrationPath, err = filepath.Abs(filepath.Dir(*configFile))

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		migrationPath = filepath.Join(migrationPath, conf.DB.MigrationPath)
	}

	if err := goose.Run(command, db, migrationPath, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}

//Grabbed this from Goose Main https://github.com/pressly/goose/blob/master/cmd/goose/main.go
func usage(flags *flag.FlagSet) func() {
	return func() {

		usagePrefix := `
db [flags] command
	`

		usageCommands := `
Commands:
  up                   	Migrate the DB to the most recent version available
  up-by-one            	Migrate the DB up by 1
  up-to VERSION        	Migrate the DB to a specific VERSION
  down                 	Roll back the version by 1
  down-to VERSION      	Roll back to a specific VERSION
  redo                 	Re-run the latest migration
  reset                	Roll back all migrations
  status               	Dump the migration status for the current DB
  version              	Print the current version of the database
  create NAME [sql|go] 	Creates new migration file with the current timestamp
  fix                  	Apply sequential ordering to migrations
	`
		fmt.Println(usagePrefix)
		fmt.Println("Flags:")
		flags.PrintDefaults()
		fmt.Println(usageCommands)
	}
}

/*
* load config from the given file into the global config struct
 */
func (conf *config) load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	jdc := json.NewDecoder(file)
	err = jdc.Decode(conf)
	if err != nil {
		return err
	}

	return nil
}

func (conf *dbConfig) postgresDBString() string {
	fmt.Printf(postgresDBString+"\n", conf.Username, conf.Password, conf.Host, conf.Port, conf.DBName, conf.SSLMode)
	return fmt.Sprintf(postgresDBString, conf.Username, conf.Password, conf.Host, conf.Port, conf.DBName, conf.SSLMode)
}
