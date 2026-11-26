package handlers

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"karsa/internal/httpx"
	"karsa/internal/middleware"
	"karsa/internal/models"
	"karsa/internal/util"
)

const (
	tokenLifetime    = 7 * 24 * time.Hour
	maxPasswordBytes = 72
)

type registerInput struct {
	Name        string `json:"name" binding:"required,notblank,max=100"`
	Email       string `json:"email" binding:"required,email,max=191"`
	Password    string `json:"password" binding:"required,min=8,max=72"`
	Role        string `json:"role" binding:"required,oneof=seeker company"`
	CompanyName string `json:"company_name" binding:"required_if=Role company,max=150"`
}

type loginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (handler *Handler) Register(ctx *gin.Context) {
	var input registerInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	companyName := strings.TrimSpace(input.CompanyName)
	if input.Role == "company" && companyName == "" {
		httpx.FieldErrors(ctx, map[string]string{"company_name": "Wajib diisi."})
		return
	}
	if len(input.Password) > maxPasswordBytes {
		httpx.FieldErrors(ctx, map[string]string{"password": "Maksimal 72."})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}

	user := models.User{
		Name:         strings.TrimSpace(input.Name),
		Email:        strings.ToLower(input.Email),
		PasswordHash: string(hash),
		Role:         input.Role,
		IsActive:     true,
	}
	if user.Role == "company" {
		user.Company = &models.Company{Name: companyName, Slug: util.UniqueSlug(companyName), Status: "pending"}
	}
	err = handler.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit("Company", "Profile").Create(&user).Error; err != nil {
			return err
		}
		if user.Company == nil {
			return nil
		}
		user.Company.UserID = user.ID
		return tx.Create(user.Company).Error
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		httpx.FieldErrors(ctx, map[string]string{"email": "Email sudah terdaftar."})
		return
	}
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	handler.respondWithToken(ctx, 201, user)
}

func (handler *Handler) Login(ctx *gin.Context) {
	var input loginInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	var user models.User
	err := handler.DB.Preload("Company").Where("email = ?", strings.ToLower(input.Email)).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		httpx.ServerError(ctx, err)
		return
	}
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		httpx.Error(ctx, 401, "Email atau password salah.")
		return
	}
	if !user.IsActive {
		httpx.Error(ctx, 403, "Akun dinonaktifkan.")
		return
	}
	handler.respondWithToken(ctx, 200, user)
}

func (handler *Handler) Me(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"user": middleware.CurrentUser(ctx)})
}

func (handler *Handler) respondWithToken(ctx *gin.Context, status int, user models.User) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(user.ID), 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(handler.Config.JWTSecret))
	if err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	ctx.JSON(status, gin.H{"user": user, "token": token})
}
