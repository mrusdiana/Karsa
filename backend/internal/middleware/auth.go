package middleware

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"karsa/internal/httpx"
	"karsa/internal/models"
)

const userKey = "currentUser"

var errInvalidToken = errors.New("token tidak valid")

func Auth(db *gorm.DB, secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := loadUser(ctx, db, secret)
		if errors.Is(err, errInvalidToken) {
			httpx.Error(ctx, 401, "Silakan masuk terlebih dahulu.")
			return
		}
		if err != nil {
			httpx.ServerError(ctx, err)
			return
		}
		if !user.IsActive {
			httpx.Error(ctx, 403, "Akun dinonaktifkan.")
			return
		}
		ctx.Set(userKey, user)
		ctx.Next()
	}
}

func OptionalAuth(db *gorm.DB, secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if user, err := loadUser(ctx, db, secret); err == nil && user.IsActive {
			ctx.Set(userKey, user)
		}
		ctx.Next()
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if CurrentUser(ctx).Role != role {
			httpx.Error(ctx, 403, "Kamu tidak punya akses.")
			return
		}
		ctx.Next()
	}
}

func RequireVerifiedCompany(ctx *gin.Context) {
	if CurrentUser(ctx).Company.Status != "verified" {
		httpx.Error(ctx, 403, "Perusahaan belum terverifikasi")
		return
	}
	ctx.Next()
}

func LimitPerMinute(limit int) gin.HandlerFunc {
	var mutex sync.Mutex
	attempts := map[string][]time.Time{}

	return func(ctx *gin.Context) {
		now := time.Now()
		expired := func(attempt time.Time) bool { return now.Sub(attempt) >= time.Minute }

		mutex.Lock()
		for address, times := range attempts {
			if times = slices.DeleteFunc(times, expired); len(times) == 0 {
				delete(attempts, address)
			} else {
				attempts[address] = times
			}
		}
		allowed := len(attempts[ctx.ClientIP()]) < limit
		if allowed {
			attempts[ctx.ClientIP()] = append(attempts[ctx.ClientIP()], now)
		}
		mutex.Unlock()

		if !allowed {
			httpx.Error(ctx, 429, "Terlalu banyak percobaan. Coba lagi dalam satu menit.")
			return
		}
		ctx.Next()
	}
}

func CurrentUser(ctx *gin.Context) *models.User {
	user, _ := ctx.Get(userKey)
	current, _ := user.(*models.User)
	return current
}

func loadUser(ctx *gin.Context, db *gorm.DB, secret string) (*models.User, error) {
	token := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(*jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithExpirationRequired())
	if err != nil {
		return nil, errInvalidToken
	}
	id, err := strconv.ParseUint(claims.Subject, 10, 63)
	if err != nil {
		return nil, errInvalidToken
	}
	var user models.User
	err = db.Preload("Company").First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errInvalidToken
	}
	return &user, err
}
