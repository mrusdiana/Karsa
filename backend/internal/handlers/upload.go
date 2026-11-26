package handlers

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"karsa/internal/httpx"
)

type uploadRule struct {
	field      string
	folder     string
	maxBytes   int64
	extensions map[string]string
	tooLarge   string
	wrongType  string
}

var (
	cvRule = uploadRule{
		field: "cv", folder: "cv", maxBytes: 2 << 20,
		extensions: map[string]string{"application/pdf": ".pdf"},
		tooLarge:   "Ukuran CV maksimal 2 MB.", wrongType: "CV harus berupa file PDF.",
	}
	logoRule = uploadRule{
		field: "logo", folder: "logos", maxBytes: 1 << 20,
		extensions: map[string]string{"image/jpeg": ".jpg", "image/png": ".png", "image/webp": ".webp"},
		tooLarge:   "Ukuran logo maksimal 1 MB.", wrongType: "Logo harus JPG, PNG, atau WebP.",
	}
)

func (handler *Handler) saveUpload(ctx *gin.Context, rule uploadRule) (string, bool) {
	reject := func(message string) (string, bool) {
		httpx.FieldErrors(ctx, map[string]string{rule.field: message})
		return "", false
	}

	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, rule.maxBytes+(1<<20))
	header, err := ctx.FormFile(rule.field)
	var sizeError *http.MaxBytesError
	if errors.As(err, &sizeError) {
		return reject(rule.tooLarge)
	}
	if err != nil {
		return reject("File wajib diisi.")
	}
	if header.Size > rule.maxBytes {
		return reject(rule.tooLarge)
	}
	file, err := header.Open()
	if err != nil {
		httpx.ServerError(ctx, err)
		return "", false
	}
	defer file.Close()

	start := make([]byte, 512)
	length, err := io.ReadFull(file, start)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		httpx.ServerError(ctx, err)
		return "", false
	}
	extension, isAllowed := rule.extensions[http.DetectContentType(start[:length])]
	if !isAllowed {
		return reject(rule.wrongType)
	}

	path := newFilePath(rule.folder, extension)
	content := io.MultiReader(bytes.NewReader(start[:length]), file)
	if err := handler.writeFile(path, content); err != nil {
		httpx.ServerError(ctx, err)
		return "", false
	}
	return path, true
}

func (handler *Handler) copyFile(source, folder string) (string, error) {
	file, err := os.Open(handler.filePath(source))
	if err != nil {
		return "", err
	}
	defer file.Close()
	path := newFilePath(folder, filepath.Ext(source))
	return path, handler.writeFile(path, file)
}

func (handler *Handler) writeFile(path string, content io.Reader) error {
	fullPath := handler.filePath(path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o750); err != nil {
		return err
	}
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o640)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(file, content)
	if err := errors.Join(copyErr, file.Close()); err != nil {
		handler.removeFile(&path)
		return err
	}
	return nil
}

func (handler *Handler) removeFile(path *string) {
	if path == nil {
		return
	}
	if err := os.Remove(handler.filePath(*path)); err != nil {
		log.Println("gagal menghapus file:", err)
	}
}

func (handler *Handler) filePath(path string) string {
	return filepath.Join(handler.Config.UploadDir, path)
}

func (handler *Handler) logoURL(path *string) *string {
	if path == nil {
		return nil
	}
	url := handler.Config.APIURL + "/uploads/" + *path
	return &url
}

func (handler *Handler) sendCV(ctx *gin.Context, path, name string) {
	fullPath := handler.filePath(path)
	if _, err := os.Stat(fullPath); err != nil {
		httpx.Error(ctx, 404, "File CV tidak ditemukan.")
		return
	}
	ctx.FileAttachment(fullPath, name)
}

func newFilePath(folder, extension string) string {
	return folder + "/" + rand.Text() + extension
}
