package imagekit

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user, password, ok := request.BasicAuth()
		file, header, err := request.FormFile("file")
		if err != nil || !ok || user != "private_abc" || password != "" {
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(writer, `{"message":"Your request does not contain private API key."}`)
			return
		}
		content, _ := io.ReadAll(file)
		if string(content) != "gambar" || header.Filename != "avatar.png" || request.FormValue("folder") != "karsa/avatars" {
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(writer, `{"message":"bad request"}`)
			return
		}
		_, _ = io.WriteString(writer, `{"url":"https://ik.imagekit.io/demo/karsa/avatars/avatar_x1.png"}`)
	}))
	defer server.Close()
	uploadEndpoint = server.URL

	url, err := Upload(context.Background(), "private_abc", "karsa/avatars", []byte("gambar"), "avatar.png")
	if err != nil || url != "https://ik.imagekit.io/demo/karsa/avatars/avatar_x1.png" {
		t.Fatalf("Upload() = %q, %v", url, err)
	}
	if _, err := Upload(context.Background(), "salah", "karsa/avatars", []byte("gambar"), "avatar.png"); err == nil {
		t.Fatal("kunci salah seharusnya gagal")
	}
}
