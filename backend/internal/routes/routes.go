package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"karsa/internal/config"
	"karsa/internal/handlers"
	"karsa/internal/httpx"
	"karsa/internal/mailer"
	"karsa/internal/middleware"
)

func New(db *gorm.DB, cfg config.Config) (*gin.Engine, error) {
	if err := httpx.SetupValidator(); err != nil {
		return nil, err
	}
	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		return nil, err
	}
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.FrontendURL},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))
	router.Static("/uploads/logos", cfg.UploadDir+"/logos")
	router.Static("/uploads/avatars", cfg.UploadDir+"/avatars")

	handler := &handlers.Handler{DB: db, Config: cfg, Mailer: mailer.Mailer{Config: cfg}}

	api := router.Group("/api")
	api.GET("/health", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"status": "ok"}) })
	api.POST("/auth/register", handler.Register)
	api.POST("/auth/login", middleware.LimitPerMinute(5), handler.Login)
	api.GET("/categories", handler.ListCategories)

	public := api.Group("", middleware.OptionalAuth(db, cfg.JWTSecret))
	public.GET("/home", handler.Home)
	public.GET("/jobs", handler.ListJobs)
	public.GET("/jobs/:id", handler.ShowJob)
	public.GET("/companies/:slug", handler.ShowCompany)

	authed := api.Group("", middleware.Auth(db, cfg.JWTSecret))
	authed.GET("/auth/me", handler.Me)
	authed.POST("/me/avatar", handler.UploadAvatar)

	seeker := authed.Group("", middleware.RequireRole("seeker"))
	seeker.POST("/jobs/:id/apply", handler.Apply)
	seeker.GET("/me/profile", handler.ShowProfile)
	seeker.PUT("/me/profile", handler.UpdateProfile)
	seeker.GET("/me/cv", handler.DownloadCV)
	seeker.POST("/me/cv", handler.UploadCV)
	seeker.GET("/me/notifications", handler.MyNotifications)
	seeker.GET("/me/applications", handler.ListMyApplications)
	seeker.GET("/me/applications/:id", handler.ShowMyApplication)
	seeker.GET("/me/applications/:id/cv", handler.DownloadMyApplicationCV)
	seeker.POST("/me/applications/:id/withdraw", handler.WithdrawApplication)
	seeker.GET("/me/saved-jobs", handler.ListSavedJobs)
	seeker.PUT("/me/saved-jobs/:id", handler.SaveJob)
	seeker.DELETE("/me/saved-jobs/:id", handler.UnsaveJob)
	seeker.GET("/me/invitations", handler.ListMyInvitations)
	seeker.POST("/me/invitations/:id/decline", handler.DeclineInvitation)

	company := authed.Group("/company", middleware.RequireRole("company"))
	company.GET("/profile", handler.ShowCompanyProfile)
	company.PUT("/profile", handler.UpdateCompanyProfile)
	company.POST("/logo", handler.UploadLogo)
	company.GET("/dashboard", handler.CompanyDashboard)
	company.GET("/jobs", handler.ListCompanyJobs)
	company.POST("/jobs", handler.CreateJob)
	company.GET("/jobs/:id", handler.ShowCompanyJob)
	company.PUT("/jobs/:id", handler.UpdateJob)
	company.DELETE("/jobs/:id", handler.DeleteJob)
	company.POST("/jobs/:id/close", handler.CloseJob)
	company.POST("/jobs/:id/open", middleware.RequireVerifiedCompany, handler.OpenJob)
	company.GET("/jobs/:id/applications", handler.ListApplicants)
	company.GET("/applications/:id", handler.ShowApplicant)
	company.PATCH("/applications/:id", handler.UpdateApplicant)
	company.GET("/applications/:id/cv", handler.DownloadApplicantCV)
	company.GET("/talents", middleware.RequireVerifiedCompany, handler.ListTalents)
	company.GET("/invitations", handler.ListSentInvitations)
	company.POST("/invitations", middleware.RequireVerifiedCompany, handler.CreateInvitation)

	admin := authed.Group("/admin", middleware.RequireRole("admin"))
	admin.GET("/dashboard", handler.AdminDashboard)
	admin.GET("/companies", handler.AdminListCompanies)
	admin.GET("/companies/:id", handler.AdminShowCompany)
	admin.POST("/companies/:id/verify", handler.VerifyCompany)
	admin.POST("/companies/:id/reject", handler.RejectCompany)
	admin.GET("/users", handler.AdminListUsers)
	admin.PATCH("/users/:id", handler.AdminUpdateUser)
	admin.GET("/jobs", handler.AdminListJobs)
	admin.POST("/jobs/:id/close", handler.AdminCloseJob)
	admin.DELETE("/jobs/:id", handler.AdminDeleteJob)
	admin.GET("/categories", handler.ListCategories)
	admin.POST("/categories", handler.CreateCategory)
	admin.PUT("/categories/:id", handler.UpdateCategory)
	admin.DELETE("/categories/:id", handler.DeleteCategory)

	return router, nil
}
