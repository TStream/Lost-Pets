package http

import (
	domain "lostpets"
	"net/http"
	"path"
	"strconv"

	"github.com/labstack/echo/v4"
)

type fileHandler struct {
	logger    domain.StructuredLogger
	router    *echo.Echo
	fileRepo  domain.FileRepo
	fileStore domain.FileStore
}

type FileInfo struct {
	ID          int
	ContentType string
}

const (
	MB = 1 << 20
)

func (h *fileHandler) initRoute(path string) {
	//File Endpoints
	h.router.GET(path+"/:id", h.handleServeFile())
	h.router.POST(path, h.handleUploadFile(path))

}

func (h *fileHandler) handleServeFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileID := c.Param("id")
		id, err := strconv.Atoi(fileID)
		if err != nil {
			return err
		}
		fileMeta, err := h.fileRepo.GetFileMeta(id)
		if err != nil {
			return err
		}

		filePath, err := h.fileStore.GetFile(fileMeta.GUID)
		if err != nil {
			return err
		}
		return c.File(filePath)
	}
}

func (h *fileHandler) handleUploadFile(location string) echo.HandlerFunc {
	type response struct {
		File FileInfo
	}
	return func(c echo.Context) error {
		// Source
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		guid, err := h.fileStore.SaveFile(src)
		if err != nil {
			return err
		}
		// Reset the read pointer
		src.Seek(0, 0)
		//detect the content-type
		buffer := make([]byte, 512)
		_, err = src.Read(buffer)
		if err != nil {
			return err
		}
		contentType := http.DetectContentType(buffer)

		fileMeta := &domain.FileMeta{
			GUID:        guid,
			ContentType: contentType,
		}
		h.fileRepo.SaveFileMeta(fileMeta)

		c.Response().Header().Set(echo.HeaderLocation, path.Join(location, strconv.Itoa(fileMeta.ID)))
		c.Response().WriteHeader(http.StatusCreated)
		return nil
	}
}
