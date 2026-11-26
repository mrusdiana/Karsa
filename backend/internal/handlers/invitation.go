package handlers

import (
	"time"

	"github.com/gin-gonic/gin"

	"karsa/internal/httpx"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

type receivedInvitation struct {
	ID        uint           `json:"id"`
	Job       jobSummary     `json:"job"`
	Company   companySummary `json:"company"`
	Message   *string        `json:"message"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `json:"expires_at"`
}

func (handler *Handler) ListMyInvitations(ctx *gin.Context) {
	query := handler.DB.Preload("Job.Company").Preload("Job.Category").Where("user_id = ?", middleware.CurrentUser(ctx).ID)
	if statuses := splitList(ctx.Query("status")); len(statuses) > 0 {
		query = query.Where("status IN ?", statuses)
	}
	var invitations []models.Invitation
	if err := query.Order("created_at DESC, id DESC").Find(&invitations).Error; err != nil {
		httpx.ServerError(ctx, err)
		return
	}
	views := make([]receivedInvitation, len(invitations))
	for index, invitation := range invitations {
		views[index] = handler.receivedInvitation(invitation)
	}
	ctx.JSON(200, views)
}

func (handler *Handler) DeclineInvitation(ctx *gin.Context) {
	id, ok := httpx.ParamID(ctx)
	if !ok {
		return
	}
	var invitation models.Invitation
	err := handler.DB.Preload("Job.Company").Preload("Job.Category").
		Where("user_id = ?", middleware.CurrentUser(ctx).ID).First(&invitation, id).Error
	if err != nil {
		httpx.FindError(ctx, err, "Undangan tidak ditemukan.")
		return
	}
	result := handler.DB.Model(&models.Invitation{}).
		Where("id = ? AND status = ?", id, "pending").Update("status", "declined")
	if result.Error != nil {
		httpx.ServerError(ctx, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		httpx.Error(ctx, 422, "Undangan sudah tidak berlaku")
		return
	}
	invitation.Status = "declined"
	ctx.JSON(200, handler.receivedInvitation(invitation))
}

func (handler *Handler) receivedInvitation(invitation models.Invitation) receivedInvitation {
	return receivedInvitation{
		ID:        invitation.ID,
		Job:       handler.publicSummary(*invitation.Job),
		Company:   handler.companySummary(*invitation.Job.Company),
		Message:   invitation.Message,
		Status:    invitation.Status,
		CreatedAt: invitation.CreatedAt,
		ExpiresAt: invitation.ExpiresAt(),
	}
}
