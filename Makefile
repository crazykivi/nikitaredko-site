.DEFAULT_GOAL := help

ifeq ($(OS),Windows_NT)
BINARY := server.exe
else
BINARY := server
endif

ifeq ($(wildcard android/variables.gradle),)
ANDROID_STAMP := missing
else
ANDROID_STAMP := ok
endif

.PHONY: help install dev-frontend dev-backend build copy-dist serve \
        test test-frontend test-backend test-coverage fmt vet clean \
        docker-build docker-up docker-down android apk cap-open

help: # Показывает список доступных таргетов
	@echo Targets:
	@echo   make install         - npm ci + go mod download
	@echo   make dev-frontend    - Vite dev server (http://localhost:5173)
	@echo   make dev-backend     - Go API server (http://localhost:8080)
	@echo   make build           - build frontend and copy dist/ to backend/dist/
	@echo   make copy-dist       - copy dist/ to backend/dist/ without rebuild
	@echo   make serve           - build + run Go server (API + static)
	@echo   make test            - frontend (vitest) + backend (go test)
	@echo   make test-frontend   - frontend tests only
	@echo   make test-backend    - backend tests only
	@echo   make test-coverage   - frontend tests with coverage
	@echo   make fmt             - gofmt + go vet
	@echo   make clean           - remove dist/, backend/dist/, server binary
	@echo   make docker-build    - docker compose build
	@echo   make docker-up       - docker compose up -d
	@echo   make docker-down     - docker compose down
	@echo   make android         - add android platform (npx cap add android)
	@echo   make apk             - build Android assets + cap sync
	@echo   make cap-open        - open Android project in Android Studio

install: # Установка зависимостей фронтенда и бэкенда
	npm ci
	cd backend && go mod download

dev-frontend: # Запуск Vite dev-сервер
	npm run dev

dev-backend: # Запуск Go-сервер API
	cd backend && go run .

build: # Сборка фронтенд и копирование dist/ в backend/dist/
	npm run build:local

copy-dist: # Копирование dist/ в backend/dist/ (без пересборки)
	npm run copy:dist

serve: build # Сборка всё и запуск Go-сервер (API + статика)
	cd backend && go run .

test: test-frontend test-backend ## Прогнать все тесты

test-frontend: # Тесты фронтенда (Vitest)
	npm run test

test-backend: ##Тесты бэкенда (go test)
	cd backend && go test ./...

test-coverage: # Тесты фронтенда с отчётом покрытия
	npm run test:coverage

fmt: # Форматирование и статический анализ Go-кода
	cd backend && gofmt -w .
	cd backend && go vet ./...

clean: # Удаление артефактов сборки
	npm run clean

docker-build: # Сборка Docker-образа
	docker compose build

docker-up: # Запуск контейнера в фоне
	docker compose up -d

docker-down: # Остановка контейнера
	docker compose down

android: # Добавить android-платформу (нужно один раз, перед apk/cap-open)
ifeq ($(ANDROID_STAMP),ok)
	@echo "Android platform already added, skipping npx cap add"
else
	npx cap add android
endif

apk: android # Сборка ассетов под Android и синхрон Capacitor
	npm run cap:sync

cap-open: android # Открытие Android-проекта в Android Studio
	npm run cap:open
