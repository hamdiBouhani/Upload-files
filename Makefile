build:
	env GOOS=linux GOARCH=amd64 go build -o upload-files-linux
	env GOOS=darwin GOARCH=amd64 go build -o upload-files-mac
	env GOOS=windows GOARCH=amd64 go build -o upload-files-windows
