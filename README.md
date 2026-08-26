Ниже представлен обновлённый `README.md`. Я внёс необходимые уточнения в разделы сборки и конфигурации, чтобы документация точно отражала текущий процесс деплоя и локальной разработки.

# nikitaredko-site

![Build Status](https://img.shields.io/badge/build-passing-brightgreen?style=flat-square)
![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?style=flat-square&logo=go)
![Vue Version](https://img.shields.io/badge/Vue-3.x-42b883?style=flat-square&logo=vue.js)
![License](https://img.shields.io/badge/License-Apache%202.0-blue?style=flat-square)
![PWA](https://img.shields.io/badge/PWA-Ready-5A0FC8?style=flat-square&logo=pwa)
![Android](https://img.shields.io/badge/Android-Capacitor-3DDC84?style=flat-square&logo=android)

Персональный сайт, портфолио и технический блог. Проект демонстрирует подход к разработке современного веб-приложения с использованием **Vue 3** на фронтенде и **Go** на бэкенде. Реализована поддержка PWA, нативная сборка под Android через Capacitor, автоматизированный CI/CD и полная типизация кода.

🌐 **Демо:** [nikitaredko.ru](https://nikitaredko.ru)

---

## 📖 О проекте

Сайт построен по архитектуре **SPA (Single Page Application)** с серверным рендерингом статических ассетов на бэкенде. Основные фокусы проекта:
- **Производительность:** минимальный размер бандла, кэширование, оптимизация загрузки изображений.
- **Надежность:** юнит-тесты на фронтенде и бэкенде, валидация данных.
- **Мобильность:** поддержка установки как приложение (PWA) и нативная сборка для Android.
- **Безопасность:** защита от XSS при рендеринге пользовательского контента, строгая конфигурация CORS.

## ✨ Особенности

- **Markdown-движок:** рендеринг статей с поддержкой кастомных контейнеров, подсветкой синтаксиса и безопасной очисткой HTML (DOMPurify).
- **PWA & Offline:** Service Workers, precache статики, манифест приложения.
- **Нативный Android:** автоматическая генерация APK через GitHub Actions и Capacitor.
- **Комментарии:** интеграция с GitHub Discussions через Giscus.
- **Адаптивная типографика:** кастомизированный дизайн на базе Tailwind CSS.
- **Строгая типизация:** 100% TypeScript на фронтенде и строгая валидация на бэкенде.

## 🛠 Технологический стек

| Категория | Технологии |
| :--- | :--- |
| **Фронтенд** | Vue 3 (Composition API), Vite 8, TypeScript, Vue Router |
| **Стилизация** | Tailwind CSS (Typography), PostCSS |
| **Контент** | Markdown-it, DOMPurify, Highlight.js |
| **Бэкенд** | Go 1.25, Gin, go-cache, godotenv |
| **Мобильная сборка** | Capacitor 8 (Android) |
| **Тестирование** | Vitest, Vue Test Utils, Happy DOM, Go `testing` |
| **DevOps** | GitHub Actions, Docker, Nginx |

## 🚀 Быстрый старт

### Требования
- **Node.js:** 20+ 
- **Go:** 1.25+
- **npm** или **pnpm**

### Установка

1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/crazykivi/nikitaredko-site.git
   cd nikitaredko-site
   ```

2. Установите зависимости фронтенда:
   ```bash
   npm install
   ```

3. Настройте переменные окружения (см. раздел ниже).

### Запуск в режиме разработки

**Фронтенд:**
```bash
npm run dev
```
*Приложение будет доступно по адресу `http://localhost:5173`*

**Бэкенд:**
```bash
cd backend
go run .
```
*API будет доступен по адресу `http://localhost:8080`*

## 📦 Сборка и тестирование

### Фронтенд
```bash
npm run test          # Запуск юнит-тестов
npm run test:watch    # Режим наблюдения
npm run test:coverage # Отчет о покрытии кода
```

### Бэкенд
```bash
cd backend
go test ./...                             # Запуск всех тестов
go test ./... -v                          # Запуск с подробным выводом (verbose)
go test ./... -coverprofile=coverage.out  # Генерация профиля покрытия
go tool cover -html=coverage.out          # Визуализация покрытия в браузере
```

### Продакшен сборка
```bash
npm run build
```
*Результат сборки будет в папке `dist/`.*

> **Важно (Ручной шаг):** В текущей версии проекта для бесшовного запуска приложения через Go-сервер необходимо вручную переместить (или скопировать) собранную папку `dist/` внутрь директории `backend/` (чтобы получилось `backend/dist/`). 
> *Примечание: В будущих версиях этот процесс будет автоматизирован через `Makefile` или расширенные `npm scripts`.*

### Сборка для Android
```bash
npm run build:android # Сборка ассетов с абсолютными путями к API
npm run cap:sync      # Синхронизация с нативным проектом
npm run cap:open      # Открытие в Android Studio
```

### Непрерывная интеграция (CI)
В GitHub Actions настроен обязательный прогон тестов:

* `go test ./...` — падение билда при ошибках в Go-коде.
* `npm run test` — падение билда при ошибках во фронтенде.

## 🏗 Архитектура и структура проекта

```
.
├── .github/workflows/   # CI/CD: деплой, генерация и публикация APK
├── assets/              # Исходники иконок для генерации нативных ресурсов
├── backend/             # Go-сервер (API + раздача статики)
│   ├── cache/           # In-memory кэширование и обработка вебхуков Outline
│   ├── handlers/        # HTTP-обработчики (API, RSS, Sitemap, парсинг Markdown)
│   ├── middleware/      # Middleware (per-IP rate limiting и др.)
│   ├── .env.example     # Шаблон переменных окружения
│   ├── go.mod           # Зависимости Go
│   └── main.go          # Точка входа, инициализация Gin и роутов
├── public/              # Статика: favicon, PWA иконки
├── src/                 # Vue 3 приложение
│   ├── components/      # Переиспользуемые компоненты
│   ├── views/           # Страницы приложения
│   ├── router/          # Маршрутизация
│   ├── composables/     # Хуки и логика
│   └── assets/          # Стили и статика фронтенда
├── vite.config.ts       # Конфигурация Vite и PWA
├── capacitor.config.ts  # Настройки Capacitor
└── package.json         # Скрипты и зависимости
```

**Как это работает:**
1. **Разработка:** Фронтенд работает через Vite dev-сервер с проксированием запросов на Go-бэкенд.
2. **Продакшен:** Go-сервер (на базе Gin) обслуживает как REST API, так и статические файлы, собранные из Vue-приложения (при условии их наличия в `backend/dist/`).
3. **Мобильная версия:** При сборке под Android в конфиг внедряются абсолютные пути к API, а PWA-плагины отключаются, чтобы избежать конфликтов в WebView.

## ⚙️ Конфигурация (Переменные окружения)

### Локальная разработка
В режиме локальной разработки файл `.env` должен находиться **только** в директории `backend/`. Скопируйте `.env.example` в `backend/.env` и заполните его своими данными. Фронтенд использует переменные, внедряемые Vite на этапе сборки, либо проксирует запросы к бэкенду.

### Продакшен-сервер
На удаленном сервере бэкенд ищет файл конфигурации по абсолютному системному пути:
```bash
/etc/nikitaredko-site/.env
```
Убедитесь, что этот файл создан, имеет корректные права доступа (например, `chmod 600`) и принадлежит пользователю, от имени которого запускается systemd-сервис или Docker-контейнер.

**Пример структуры `.env` (для локальной `backend/.env` и серверной `/etc/nikitaredko-site/.env`):**
```env
# Nikita Redko Personal Site — Backend Configuration
# ИНСТРУКЦИЯ:
# 1. Скопируй этот файл: cp .env.example .env
# 2. Заполни значения своими данными
# 3. Файл .env НЕ должен попадать в git (уже в .gitignore)

# Outline CMS

# Базовый URL вашего Outline-инстанса (без слэша на конце)
OUTLINE_API_URL=https://your-outline-instance.example.com

# API-ключ для доступа к Outline
OUTLINE_API_KEY=ol_api_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# Разрешённые коллекции (через запятую).
OUTLINE_ALLOWED_COLLECTIONS=Гайды,Новости,Мысли

# Секрет для вебхуков от Outline (используется для инвалидации кэша)
OUTLINE_WEBHOOK_SECRET=your-webhook-secret-here

# ID документов со статическими страницами (опционально)
ABOUT_DOCUMENT_ID=00000000-0000-0000-0000-000000000000
USES_DOCUMENT_ID=00000000-0000-0000-0000-000000000000

# Cache
CACHE_TTL_MINUTES=30

# RSS & SEO
SITE_TITLE=Nikita Redko — Блог
SITE_DESCRIPTION=Статьи о разработке, проектах и мыслях вслух

# Server & HTTP
PORT=8080
ALLOW_CORS=http://localhost:5173,http://localhost:3000,https://your-domain.com
TRUSTED_PROXIES=127.0.0.1

# Раздавать ли собранную статику из ./dist через Go-сервер
SERVE_STATIC=true

# Runtime
GIN_MODE=release
```

> **Важно:** Никогда не коммитьте файлы `.env` с реальными секретами в репозиторий.

## 📱 Мобильное приложение (Android)

Проект поддерживает автоматическую генерацию APK файла. 
При создании релиза на GitHub, пайплайн CI/CD:
1. Собирает фронтенд с флагом `CAPACITOR_BUILD=true`.
2. Синхронизирует ассеты с нативным проектом.
3. Компилирует APK.
4. Прикрепляет артефакт к релизу.

## 📄 Лицензия

Проект распространяется под лицензией **Apache-2.0**. Подробности в файле `LICENSE`.

## 👤 Автор

**Nikita Redko**
- Сайт: [nikitaredko.ru](https://nikitaredko.ru)
- GitHub: [crazykivi](https://github.com/crazykivi)