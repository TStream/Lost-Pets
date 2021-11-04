package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	domain "lostpets"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	Config struct {
		Host         string      `json:"host"`
		Port         int         `json:"port"`
		FileLocation string      `json:"fileLocation"`
		Debug        bool        `json:"debug"`
		Email        EmailConfig `json:"emailSettings"`
	}

	EmailConfig struct {
		Host     string           `json:"host"`
		Port     int              `json:"port"`
		User     string           `json:"user"`
		Password string           `json:"password"`
		LinkBase string           `json:"linkBase"`
		Template TemplateSettings `json:"template"`
	}

	TemplateSettings struct {
		FilePath string `json:"filePath"`
		Default  string `json:"default"`
	}

	response struct {
		Data interface{} `json:"data"`
		Meta interface{} `json:"meta,omitempty"`
	}

	emailer struct {
		config EmailConfig
		logger domain.StructuredLogger
	}

	apiPetType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

const (

	// url params
	paramShipmentID = "shipment_id"
	paramShareLink  = "share_link"

	//url paths
	filePath      = "/pet-pictures"
	postingsPath  = "/postings"
	sightingsPath = "/sightings"
)

/*StartServer configures and starts a new http server*/
func StartServer(config Config, db domain.LostPetsRepo, fileDb domain.FileRepo, fileStore domain.FileStore, logger domain.StructuredLogger, version string) {

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderLocation},
	}))

	emailer := emailer{config: config.Email, logger: logger}

	fileHandler := fileHandler{logger: logger, fileRepo: fileDb, fileStore: fileStore, router: e}
	fileHandler.initRoute(filePath)

	postingHandler := postingsHandler{logger: logger, router: e, repo: db, emailer: emailer}
	postingHandler.initRoute(postingsPath)

	sightingHandler := sightingsHandler{logger: logger, router: e, repo: db, emailer: emailer}
	sightingHandler.initRoute(sightingsPath)

	e.GET("/pet-types", getPetTypesHandler(db))

	e.GET("/", apiInfoHandler(version, config.Debug))

	if config.Debug {
		data, _ := json.MarshalIndent(e.Routes(), "", "  ")
		log.Print(string(data))
	}

	if err := e.Start(":8080"); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func apiInfoHandler(version string, debug bool) echo.HandlerFunc {
	fmt.Println("version", version, "debug", debug)

	type apiInfo struct {
		Version string `json:"version"`
	}

	if len(version) == 0 {
		version = "development"
	}

	if debug {
		version += "-debug"
	}

	return func(c echo.Context) error {
		result := response{Data: apiInfo{Version: version}}
		return c.JSON(http.StatusOK, result)
	}
}

func getPetTypesHandler(db domain.LostPetsRepo) echo.HandlerFunc {

	return func(c echo.Context) error {
		types, err := db.GetPetTypes()
		if err != nil {
			return err
		}
		apiTypes := []apiPetType{}
		for _, t := range types {
			apiTypes = append(apiTypes, apiPetType{ID: t.ID, Name: t.Name})
		}
		return c.JSON(http.StatusOK, apiTypes)
	}

}

func (e emailer) emailMatches(mType, toEmail, guid string) error {
	t, err := template.ParseFiles(e.config.Template.FilePath)
	if err != nil {
		e.logger.Error(err.Error())
	}
	data := struct {
		Type string
		Link string
	}{
		mType,
		e.config.LinkBase + guid,
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", e.config.User, e.config.Password, e.config.Host)

	if err := smtp.SendMail(e.config.Host+":"+strconv.Itoa(e.config.Port), auth, e.config.User, []string{toEmail}, buf.Bytes()); err != nil {
		return err
	}

	return nil
}

//custom time type to unmarshal time formats
type Datetime struct {
	time.Time
}

func (t *Datetime) UnmarshalJSON(input []byte) error {
	strInput := strings.Trim(string(input), `"`)
	newTime, err := time.Parse("2006-01-02T15:04:05.000Z", strInput)
	if err != nil {
		return err
	}

	t.Time = newTime
	return nil
}
