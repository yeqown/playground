docker:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./main_k8s.go
	docker build -t yeqown/simple-server:v1.0.0 -f Dockerfile .