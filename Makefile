GO=${GOROOT}/bin/go

local-build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux ${GO} build -ldflags '-extldflags "-static"' -o dist/bot_amd64 cmd/tg-bot-debug/tg-bot-service.go
