default: fmt build test vet

fmt:
	go fmt ./...

vet:
	go vet ./...

build:
	@echo "🐳"
	docker build -t xapp -f docker/Dockerfile .

xapp:
	go build -o bin/xapp cmd/xapp/main.go

test:
	go test ./...

release:
	mkdir -p dist
	GO111MODULE=on go build -o dist/xapp-darwin-amd64
	GO111MODULE=on;GOOS=linux;go build -o dist/xapp-linux-amd64
	GO111MODULE=on;GOOS=windows;go build -o dist/xapp-windows-amd64
	GO111MODULE=on;GOOS=linux GOARCH=arm GOARM=5;go build -o dist/xapp-linux-arm5

deploy: build
	docker stack deploy -c docker/docker-compose-test.yml demo

run: build
	./dist/xapp-darwin-amd64

mock:
	mockgen -source internal/app/imageprocessor/imageprocessor.go -destination internal/app/imageprocessor/mock_imageprocessor/mock_imageprocessor.go -package mock_imageprocessor
	mockgen -source internal/app/filehandler/filehandler.go -destination internal/app/filehandler/mock_filehandler/mock_filehandler.go -package mock_filehandler
	mockgen -source internal/app/persistence/storage.go -destination internal/app/persistence/mock_storage/mock_storage.go -package mock_storage
	mockgen -source internal/app/imageloader/imageloader.go -destination internal/app/imageloader/mock_imageloader/mock_imageloader.go -package mock_imageloader

localdb:
	docker run \
        	-e POSTGRES_USER=xapp \
        	-e POSTGRES_PASSWORD=xapp123 \
        	-e POSTGRES_DB=xapp \
        	-v xappdb:/var/lib/postgresql/data \
        	-p 5432:5432 \
        	postgres:10.6