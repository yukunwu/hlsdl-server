tag = hlsdl-server:1.0
os = linux
arch = amd64

gbuild:
	GOOS=$(os) GOARCH=$(arch) go build -v -o bin/hlsdl_server
dbuild:
	docker build -f Dockerfile -t $(tag) .
clear:
	rm -rf bin/*
deploy: gbuild dbuild
	docker-compose up -d