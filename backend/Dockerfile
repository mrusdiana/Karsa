FROM golang:1.27-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/seed ./cmd/seed

FROM alpine:3.22
RUN addgroup -S karsa && adduser -S karsa -G karsa \
 && mkdir -p /app/uploads && chown karsa:karsa /app/uploads
WORKDIR /app
COPY --from=build /out/server /out/seed ./
COPY migrations ./migrations
USER karsa
ENV PORT=8080 UPLOAD_DIR=/app/uploads TZ=Asia/Jakarta
EXPOSE 8080
CMD ["./server"]
