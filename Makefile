NAME=s3cat

${NAME}:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags '-w' -o $(NAME) main.go
