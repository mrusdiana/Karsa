package handlers

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"karsa/internal/httpx"
	"karsa/internal/middleware"
	"karsa/internal/models"
	"karsa/internal/util"
)

const profileFieldCount = 6

type profileInput struct {
	Headline        string   `json:"headline" binding:"max=150"`
	Phone           string   `json:"phone" binding:"max=20"`
	Location        string   `json:"location" binding:"max=100"`
	Education       string   `json:"education" binding:"omitempty,oneof=SMA/SMK D3 S1 S2"`
	Major           string   `json:"major" binding:"max=100"`
	GraduationYear  *int     `json:"graduation_year" binding:"omitempty,min=1970,max=2100"`
	ExperienceYears int      `json:"experience_years" binding:"min=0,max=30"`
	Skills          []string `json:"skills" binding:"max=15,dive,max=40"`
	Bio             string   `json:"bio" binding:"max=2000"`
	PortfolioURL    string   `json:"portfolio_url" binding:"omitempty,http_url,max=255"`
	IsDiscoverable  bool     `json:"is_discoverable"`
}

type cvView struct {
	Name       string    `json:"name"`
	Size       int64     `json:"size"`
	UploadedAt time.Time `json:"uploaded_at"`
	URL        string    `json:"url"`
}

type completeness struct {
	Percent int      `json:"percent"`
	Missing []string `json:"missing"`
}

type publicProfile struct {
	Headline        *string  `json:"headline"`
	Location        *string  `json:"location"`
	Education       *string  `json:"education"`
	Major           *string  `json:"major"`
	GraduationYear  *int     `json:"graduation_year"`
	ExperienceYears int      `json:"experience_years"`
	Skills          []string `json:"skills"`
}

type contactProfile struct {
	publicProfile
	Phone        *string `json:"phone"`
	Bio          *string `json:"bio"`
	PortfolioURL *string `json:"portfolio_url"`
}

type profileView struct {
	contactProfile
	CV             *cvView      `json:"cv"`
	IsDiscoverable bool         `json:"is_discoverable"`
	Completeness   completeness `json:"completeness"`
}

func (handler *Handler) ShowProfile(ctx *gin.Context) {
	if profile, ok := handler.currentProfile(ctx); ok {
		handler.respondProfile(ctx, profile)
	}
}

func (handler *Handler) UpdateProfile(ctx *gin.Context) {
	var input profileInput
	if !httpx.BindJSON(ctx, &input) {
		return
	}
	profile, ok := handler.currentProfile(ctx)
	if !ok {
		return
	}
	profile.Headline = optional(input.Headline)
	profile.Phone = optional(input.Phone)
	profile.Location = optional(input.Location)
	profile.Education = optional(input.Education)
	profile.Major = optional(input.Major)
	profile.GraduationYear = input.GraduationYear
	profile.ExperienceYears = input.ExperienceYears
	profile.Skills = normalizeSkills(input.Skills)
	profile.Bio = optional(input.Bio)
	profile.PortfolioURL = optional(input.PortfolioURL)

	isComplete := len(profile.Missing()) == 0
	if input.IsDiscoverable && !profile.IsDiscoverable && !isComplete {
		httpx.FieldErrors(ctx, map[string]string{"is_discoverable": "Lengkapi profil dulu sebelum mengaktifkan ini."})
		return
	}
	profile.IsDiscoverable = input.IsDiscoverable && isComplete
	if err := handler.DB.Omit("cv_path").Save(&profile).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	handler.respondProfile(ctx, profile)
}

func (handler *Handler) UploadCV(ctx *gin.Context) {
	profile, ok := handler.currentProfile(ctx)
	if !ok {
		return
	}
	path, ok := handler.saveUpload(ctx, cvRule)
	if !ok {
		return
	}
	previous := profile.CVPath
	if err := handler.DB.Model(&profile).Update("cv_path", path).Error; err != nil {
		handler.removeFile(&path)
		httpx.ServerError(ctx, err)
		return
	}
	handler.removeFile(previous)
	handler.respondProfile(ctx, profile)
}

func (handler *Handler) DownloadCV(ctx *gin.Context) {
	profile, ok := handler.currentProfile(ctx)
	if !ok {
		return
	}
	if profile.CVPath == nil {
		httpx.Error(ctx, 404, "File CV tidak ditemukan.")
		return
	}
	handler.sendCV(ctx, *profile.CVPath, cvName(middleware.CurrentUser(ctx).Name))
}

func (handler *Handler) seekerProfile(userID uint) (models.SeekerProfile, error) {
	profile := models.SeekerProfile{UserID: userID}
	err := handler.DB.Where("user_id = ?", userID).Attrs(models.SeekerProfile{Skills: []string{}}).FirstOrCreate(&profile).Error
	return profile, err
}

func (handler *Handler) currentProfile(ctx *gin.Context) (models.SeekerProfile, bool) {
	profile, err := handler.seekerProfile(middleware.CurrentUser(ctx).ID)
	if err != nil {
		httpx.ServerError(ctx, err)
	}
	return profile, err == nil
}

func (handler *Handler) respondProfile(ctx *gin.Context, profile models.SeekerProfile) {
	missing := profile.Missing()
	view := profileView{
		contactProfile: newContactProfile(profile),
		IsDiscoverable: profile.IsDiscoverable,
		Completeness:   completeness{Percent: (profileFieldCount - len(missing)) * 100 / profileFieldCount, Missing: missing},
	}
	if profile.CVPath != nil {
		if info, err := os.Stat(handler.filePath(*profile.CVPath)); err == nil {
			view.CV = &cvView{
				Name:       cvName(middleware.CurrentUser(ctx).Name),
				Size:       info.Size(),
				UploadedAt: info.ModTime(),
				URL:        handler.Config.APIURL + "/api/me/cv",
			}
		}
	}
	ctx.JSON(200, view)
}

func newPublicProfile(profile models.SeekerProfile) publicProfile {
	return publicProfile{
		Headline:        profile.Headline,
		Location:        profile.Location,
		Education:       profile.Education,
		Major:           profile.Major,
		GraduationYear:  profile.GraduationYear,
		ExperienceYears: profile.ExperienceYears,
		Skills:          append([]string{}, profile.Skills...),
	}
}

func newContactProfile(profile models.SeekerProfile) contactProfile {
	return contactProfile{
		publicProfile: newPublicProfile(profile),
		Phone:         profile.Phone,
		Bio:           profile.Bio,
		PortfolioURL:  profile.PortfolioURL,
	}
}

func cvName(userName string) string {
	return "CV-" + util.Slugify(userName) + ".pdf"
}
