.PHONY: all build up run down clean

# Image names
DOWNLOADER_IMAGE = kalandar5862/video-downloader
AUDIO_EXTRACTOR_IMAGE = audio_extractor
TRANSCRIBER_IMAGE = transcriber
MAIN_SERVICE_IMAGE = main_service

# Container names
DOWNLOADER_CONTAINER = downloader_container
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

# Start all containers
up: create-network
	# Start Audio Extractor with videos folder mounted
	docker run -d --name $(AUDIO_EXTRACTOR_CONTAINER) --network $(NETWORK) -p 5001:5001 \
		-v $(PWD)/results/videos:/app/results/videos \
		$(AUDIO_EXTRACTOR_IMAGE)

	# Wait for Audio Extractor to start
	@echo "Waiting for Audio Extractor to start..."
	@until curl -s http://localhost:5001/extract >/dev/null; do sleep 1; done

	# Start Transcriber with videos and results folders mounted
	docker run -d --name $(TRANSCRIBER_CONTAINER) --network $(NETWORK) -p 5002:5002 \
		-v $(PWD)/results/videos:/app/results/videos \
		-v $(PWD)/results:/app/results \
		$(TRANSCRIBER_IMAGE)

	# Wait for Transcriber to start
	@echo "Waiting for Transcriber to start..."
	@until curl -s http://localhost:5002/transcribe >/dev/null; do sleep 1; done

	# Start Downloader with videos and results folders mounted
	docker run -d --name $(DOWNLOADER_CONTAINER) --network $(NETWORK) -p 5011:5011 \
		-v $(PWD)/results/videos:/app/results/videos \
		-v $(PWD)/results:/app/results \
		$(DOWNLOADER_IMAGE)

	# Wait for Downloader to start
	@echo "Waiting for Downloader to start..."
	@until curl -s http://localhost:5011/download >/dev/null; do sleep 1; done

	# Start Main Service with videos and results folders mounted
	docker run -d --name $(MAIN_SERVICE_CONTAINER) --network $(NETWORK) -p 8080:8080 \
		-v $(PWD)/results/videos:/app/results/videos \
		-v $(PWD)/results:/app/results \
		$(MAIN_SERVICE_IMAGE)

# Start all containers after building images
run: build up

# Stop and remove all containers
down:
	@docker stop $(AUDIO_EXTRACTOR_CONTAINER) $(TRANSCRIBER_CONTAINER) $(MAIN_SERVICE_CONTAINER) $(DOWNLOADER_CONTAINER) || true
	@docker rm -f $(AUDIO_EXTRACTOR_CONTAINER) $(TRANSCRIBER_CONTAINER) $(MAIN_SERVICE_CONTAINER) $(DOWNLOADER_CONTAINER) || true

# Clean up containers and images
clean: down
	@docker network rm $(NETWORK) || true
	@docker rmi -f $(AUDIO_EXTRACTOR_IMAGE) $(TRANSCRIBER_IMAGE) $(MAIN_SERVICE_IMAGE) $(DOWNLOADER_CONTAINER) || true

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