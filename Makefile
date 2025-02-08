include config.env
export

.PHONY: all build up run down clean

# Image names
VIDEO_DOWNLOADER_IMAGE = kalandar5862/video-downloader
AUDIO_EXTRACTOR_IMAGE = audio_extractor
TRANSCRIBER_IMAGE = transcriber
MAIN_SERVICE_IMAGE = main_service

# Container names
VIDEO_DOWNLOADER_CONTAINER = downloader_container
AUDIO_EXTRACTOR_CONTAINER = audio_extractor_container
TRANSCRIBER_CONTAINER = transcriber_container
MAIN_SERVICE_CONTAINER = main_service_container

# Custom network name
NETWORK = transcriber_network

# Build all images
build: build-audio-extractor build-transcriber build-main-service

build-audio-extractor:
	docker build -t $(AUDIO_EXTRACTOR_IMAGE) ./audio_extractor

build-transcriber:
	docker build -t $(TRANSCRIBER_IMAGE) ./transcriber

build-main-service:
	docker build -t $(MAIN_SERVICE_IMAGE) ./main_service

# Create network if it doesn't exist
create-network:
	@docker network inspect $(NETWORK) > /dev/null 2>&1 || \
		docker network create $(NETWORK)

# Start Audio Extractor with videos folder mounted
up-audio-extractor:
	@docker rm -f $(AUDIO_EXTRACTOR_CONTAINER) || true
	docker run -d --name $(AUDIO_EXTRACTOR_CONTAINER) --network $(NETWORK) --env-file ./config.env -p $(AUDIO_EXTRACTOR_PORT):$(AUDIO_EXTRACTOR_PORT) \
		-v $(PWD)/results/videos:/app/results/videos \
		$(AUDIO_EXTRACTOR_IMAGE)

	@echo "Waiting for Audio Extractor to start..."
	@until curl -s http://localhost:5001/extract >/dev/null; do sleep 1; done

# Start Transcriber with videos and results folders mounted
up-transcriber:
	@docker rm -f $(TRANSCRIBER_CONTAINER) || true
	docker run -d --name $(TRANSCRIBER_CONTAINER) --network $(NETWORK) --env-file ./config.env -p $(TRANSCRIBER_PORT):$(TRANSCRIBER_PORT) \
		-v $(PWD)/results/videos:/app/results/videos \
		-v $(PWD)/results:/app/results \
		$(TRANSCRIBER_IMAGE)

	@echo "Waiting for Transcriber to start..."
	@until curl -s http://localhost:5002/transcribe >/dev/null; do sleep 1; done

# Start Downloader with videos and results folders mounted
up-video-downloader:
	@docker rm -f $(VIDEO_DOWNLOADER_CONTAINER) || true
	docker run -d --name $(VIDEO_DOWNLOADER_CONTAINER) --network $(NETWORK) --env-file ./config.env -p $(VIDEO_DOWNLOADER_PORT):5011 \
		-v $(PWD)/results/videos:/app/results/videos \
		-v $(PWD)/results:/app/results \
		$(VIDEO_DOWNLOADER_IMAGE)

	@echo "Waiting for Downloader to start..."
	@until curl -s http://localhost:5011/download >/dev/null; do sleep 1; done

# Start Main Service with videos and results folders mounted
up-main-service:
	@docker rm -f $(MAIN_SERVICE_CONTAINER) || true
	docker run -d --name $(MAIN_SERVICE_CONTAINER) --network $(NETWORK) --env-file ./config.env -p $(MAIN_SERVICE_PORT):8080 \
		-v $(PWD)/results/videos:/app/results/videos \
		-v $(PWD)/results:/app/results \
		$(MAIN_SERVICE_IMAGE)

	@echo "Waiting for Main service to start..."
	@until curl -s http://localhost:$(MAIN_SERVICE_PORT)/download >/dev/null; do sleep 1; done

# Start all containers
up: create-network up-audio-extractor up-transcriber up-video-downloader up-main-service

# Start all containers after building images
run: build up

# Stop and remove all containers
down:
	@docker stop $(AUDIO_EXTRACTOR_CONTAINER) $(TRANSCRIBER_CONTAINER) $(MAIN_SERVICE_CONTAINER) $(VIDEO_DOWNLOADER_CONTAINER) || true
	@docker rm -f $(AUDIO_EXTRACTOR_CONTAINER) $(TRANSCRIBER_CONTAINER) $(MAIN_SERVICE_CONTAINER) $(VIDEO_DOWNLOADER_CONTAINER) || true

# Clean up containers and images
clean: down
	@docker network rm $(NETWORK) || true
	@docker rmi -f $(AUDIO_EXTRACTOR_IMAGE) $(TRANSCRIBER_IMAGE) $(MAIN_SERVICE_IMAGE) $(VIDEO_DOWNLOADER_CONTAINER) || true

rerun-main:
	docker stop $(MAIN_SERVICE_CONTAINER)
	docker rm -f $(MAIN_SERVICE_CONTAINER)
	@docker rmi -f $(MAIN_SERVICE_IMAGE) || true
	docker build -t $(MAIN_SERVICE_IMAGE) ./main_service
	# Start Main Service with videos and results folders mounted
	docker run -d --name $(MAIN_SERVICE_CONTAINER) --network $(NETWORK) -p 8080:8080 \
		-v $(PWD)/results/videos:/app/results/videos \
		-v $(PWD)/results:/app/results \
		$(MAIN_SERVICE_IMAGE)

lint: lint-audio-extractor lint-main-service lint-transcriber

lint-audio-extractor:
	golangci-lint run ./audio_extractor/*.go

lint-main-service:
	golangci-lint run ./main_service/...

lint-transcriber:
	golangci-lint run ./transcriber/*.go