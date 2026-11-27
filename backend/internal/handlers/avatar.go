package handlers

import (
	"bytes"
	"strings"

	"github.com/gin-gonic/gin"

	"karsa/internal/httpx"
	"karsa/internal/imagekit"
	"karsa/internal/middleware"
	"karsa/internal/models"
)

func (handler *Handler) UploadAvatar(ctx *gin.Context) {
	user := middleware.CurrentUser(ctx)
	content, extension, ok := readUpload(ctx, avatarRule)
	if !ok {
		return
	}

	var url string
	localPath := ""
	if handler.Config.ImageKitKey != "" {
		uploaded, err := imagekit.Upload(ctx.Request.Context(), handler.Config.ImageKitKey, "karsa/avatars", content, "avatar"+extension)
		if err != nil {
			httpx.Error(ctx, 502, "Gagal mengunggah foto. Coba lagi.")
			return
		}
		url = uploaded
	} else {
		localPath = newFilePath(avatarRule.folder, extension)
		if err := handler.writeFile(localPath, bytes.NewReader(content)); err != nil {
			httpx.ServerError(ctx, err)
			return
		}
		url = handler.Config.APIURL + "/uploads/" + localPath
	}

	if err := handler.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("avatar_url", url).Error; err != nil {
		if localPath != "" {
			handler.removeFile(&localPath)
		}
		httpx.ServerError(ctx, err)
		return
	}
	if previous := handler.localAvatarPath(user.AvatarURL); previous != nil {
		handler.removeFile(previous)
	}
	user.AvatarURL = &url
	ctx.JSON(200, gin.H{"user": user})
}

func (handler *Handler) localAvatarPath(url *string) *string {
	if url == nil {
		return nil
	}
	path, isLocal := strings.CutPrefix(*url, handler.Config.APIURL+"/uploads/")
	if !isLocal || !strings.HasPrefix(path, avatarRule.folder+"/") || strings.Contains(path, "..") {
		return nil
	}
	return &path
}
