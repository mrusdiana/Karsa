package httpx

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"gorm.io/gorm"
)

func SetupValidator() error {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return errors.New("mesin validator tidak dikenali")
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
	return validate.RegisterValidation("notblank", validators.NotBlank)
}

func Error(ctx *gin.Context, status int, message string) {
	ctx.AbortWithStatusJSON(status, gin.H{"message": message})
}

func ServerError(ctx *gin.Context, err error) {
	log.Printf("%s %s: %v", ctx.Request.Method, ctx.Request.URL.Path, err)
	Error(ctx, 500, "Terjadi kesalahan. Coba lagi.")
}

func FindError(ctx *gin.Context, err error, message string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Error(ctx, 404, message)
		return
	}
	ServerError(ctx, err)
}

func FieldErrors(ctx *gin.Context, fields map[string]string) {
	ctx.AbortWithStatusJSON(422, gin.H{"message": "Data tidak valid.", "errors": fields})
}

func BindJSON(ctx *gin.Context, input any) bool {
	err := ctx.ShouldBindJSON(input)
	if err == nil {
		return true
	}
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		Error(ctx, 422, "Format data tidak valid.")
		return false
	}
	fields := map[string]string{}
	for _, fieldError := range validationErrors {
		_, path, _ := strings.Cut(fieldError.Namespace(), ".")
		fields[path] = message(fieldError)
	}
	FieldErrors(ctx, fields)
	return false
}

func message(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required", "required_if", "notblank":
		return "Wajib diisi."
	case "email":
		return "Format email tidak valid."
	case "http_url":
		return "Tautan harus diawali http:// atau https://."
	case "min":
		return "Minimal " + fieldError.Param() + "."
	case "max":
		return "Maksimal " + fieldError.Param() + "."
	case "oneof":
		return "Pilihan tidak valid."
	case "datetime":
		return "Format tanggal harus YYYY-MM-DD."
	}
	return "Tidak valid."
}

func ParamID(ctx *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 63)
	if err != nil {
		Error(ctx, 404, "Data tidak ditemukan.")
		return 0, false
	}
	return uint(id), true
}
