# Makefile

.PHONY: all build run clean

# Имена образов
AUDIO_EXTRACTOR_IMAGE = audio_extractor
TRANSCRIBER_IMAGE = transcriber
MAIN_SERVICE_IMAGE = main_service

# Имена контейнеров
AUDIO_EXTRACTOR_CONTAINER = audio_extractor_container
TRANSCRIBER_CONTAINER = transcriber_container
MAIN_SERVICE_CONTAINER = main_service_container

# Имя пользовательской сети
NETWORK = transcriber_network

# Сборка всех образов
build: build-audio-extractor build-transcriber build-main-service

build-audio-extractor:
	docker build -t $(AUDIO_EXTRACTOR_IMAGE) ./audio_extractor

build-transcriber:
	docker build -t $(TRANSCRIBER_IMAGE) ./transcriber

build-main-service:
	docker build -t $(MAIN_SERVICE_IMAGE) ./main_service

# Создание сети
create-network:
	docker network create $(NETWORK) || true

# Запуск всех контейнеров
run: build create-network
	# Запускаем Audio Extractor с монтированием папки videos
	docker run -d --name $(AUDIO_EXTRACTOR_CONTAINER) --network $(NETWORK) -p 5001:5001 \
		-v $(PWD)/videos:/app/videos \
		$(AUDIO_EXTRACTOR_IMAGE)

	# Ждем пока Audio Extractor запустится
	@echo "Ожидание запуска Audio Extractor..."
	@until curl -s http://localhost:5001/extract >/dev/null; do sleep 1; done

	# Запускаем Transcriber с монтированием папок videos и results
	docker run -d --name $(TRANSCRIBER_CONTAINER) --network $(NETWORK) -p 5002:5002 \
		-v $(PWD)/videos:/app/videos \
		-v $(PWD)/results:/app/results \
		$(TRANSCRIBER_IMAGE)

	# Ждем пока Transcriber запустится
	@echo "Ожидание запуска Transcriber..."
	@until curl -s http://localhost:5002/transcribe >/dev/null; do sleep 1; done

	# Запускаем Main Service с монтированием папок videos и results
	docker run --name $(MAIN_SERVICE_CONTAINER) --network $(NETWORK) \
		-v $(PWD)/videos:/app/videos \
		-v $(PWD)/results:/app/results \
		$(MAIN_SERVICE_IMAGE)

# Очистка контейнеров и образов
clean:
	docker rm -f $(AUDIO_EXTRACTOR_CONTAINER) || true
	docker rm -f $(TRANSCRIBER_CONTAINER) || true
	docker rm -f $(MAIN_SERVICE_CONTAINER) || true
	docker network rm $(NETWORK) || true
	docker rmi -f $(AUDIO_EXTRACTOR_IMAGE) $(TRANSCRIBER_IMAGE) $(MAIN_SERVICE_IMAGE) || true
