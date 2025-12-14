# Инструкция по установке и тестированию TeleMonitor AI

## ✅ Что уже готово

- ✅ Код загружен в GitHub: https://github.com/LuykLembit/asist_app.git
- ✅ База данных полностью спроектирована
- ✅ Создан `main.go` для тестирования фундамента
- ✅ Docker конфигурация готова

## 📦 Требуемое ПО

### Вариант 1: Локальная разработка

1. **Go 1.22+**
   - Скачать: https://go.dev/dl/
   - Установить и добавить в PATH
   - Проверить: `go version`

2. **PostgreSQL 16**
   - Скачать: https://www.postgresql.org/download/windows/
   - Или использовать Docker (см. Вариант 2)

### Вариант 2: Docker (Рекомендуется)

1. **Docker Desktop**
   - Скачать: https://www.docker.com/products/docker-desktop/
   - Установить и запустить
   - Проверить: `docker --version` и `docker-compose --version`

## 🔧 Настройка проекта

### Шаг 1: Получить учетные данные

#### 1.1 Telegram API (для Userbot)
1. Перейти на https://my.telegram.org
2. Войти с вашим номером телефона
3. Перейти в "API development tools"
4. Создать приложение
5. Сохранить `app_id` и `app_hash`

#### 1.2 Telegram Bot Token
1. Открыть @BotFather в Telegram
2. Отправить `/newbot`
3. Следовать инструкциям
4. Сохранить токен (формат: `1234567890:ABCdefGHIjklMNOpqrsTUVwxyz`)

#### 1.3 Ваш Telegram User ID
1. Открыть @userinfobot в Telegram
2. Отправить любое сообщение
3. Сохранить ваш ID (число, например: `123456789`)

#### 1.4 ZhipuAI API Key
1. Зарегистрироваться на https://open.bigmodel.cn/
2. Перейти в раздел API Keys
3. Создать новый ключ
4. Сохранить API key

### Шаг 2: Создать .env файл

```bash
cp .env.example .env
```

Отредактировать `.env` и заполнить реальными данными:

```env
# Telegram Userbot Credentials
TG_APP_ID=12345678              # <- Ваш app_id
TG_APP_HASH=abcdef1234567890    # <- Ваш app_hash

# Telegram Bot Credentials
TG_BOT_TOKEN=1234567890:ABCdef  # <- Ваш bot token

# Admin Access
TG_ADMIN_ID=123456789           # <- Ваш user ID

# ZhipuAI API Configuration
ZHIPU_API_KEY=your_actual_key   # <- Ваш API key
```

## 🚀 Запуск тестирования

### Вариант A: Docker (Проще)

```bash
# Запустить базу данных и приложение
docker-compose up -d

# Посмотреть логи
docker-compose logs -f telemonitor

# Остановить
docker-compose down
```

**Ожидаемый вывод:**
```
🚀 TeleMonitor AI - Starting...
📋 Loading configuration...
✅ Configuration loaded successfully
🔍 Validating configuration...
✅ Configuration is valid
🗄️  Connecting to PostgreSQL...
✅ Database connected successfully
✅ Database ping successful
🔄 Running database migrations...
✅ Database migrations completed successfully
🧪 Testing repository layer...
  ✓ SessionRepository initialized
  ✓ MonitoredChatRepository initialized (chats: 0)
  ✓ RawMessageRepository initialized
  ✓ TriggerRepository initialized (triggers: 0)
  ✓ DailyReportRepository initialized
✅ Repository layer is functional
🎉 System initialization complete!
```

### Вариант B: Локальная разработка

```bash
# Установить зависимости
go mod download

# Запустить PostgreSQL отдельно
docker-compose up -d postgres

# Настроить подключение к локальной БД в .env
# (раскомментировать DB_* переменные)

# Собрать приложение
go build -o telemonitor.exe ./cmd/telemonitor

# Запустить
./telemonitor.exe
```

## 🧪 Проверка результатов

### 1. Проверка базы данных

```bash
# Подключиться к PostgreSQL
docker exec -it telemonitor_db psql -U telemonitor_user -d telemonitor

# Посмотреть таблицы
\dt

# Должны быть созданы:
# - session_storage
# - monitored_chats
# - raw_messages
# - triggers
# - daily_reports

# Выйти
\q
```

### 2. Проверка структуры таблиц

```sql
-- Посмотреть структуру таблицы
\d monitored_chats

-- Проверить индексы
\di
```

## ❌ Возможные проблемы

### Ошибка: "Failed to load configuration"
- Проверьте, что файл `.env` создан
- Убедитесь, что все обязательные поля заполнены

### Ошибка: "Database connection failed"
- Проверьте, что PostgreSQL запущен: `docker-compose ps`
- Проверьте логи БД: `docker-compose logs postgres`

### Ошибка: "Configuration validation failed: TG_APP_ID is not set"
- В `.env` должны быть числовые значения без кавычек:
  ```
  TG_APP_ID=12345678
  ```
  НЕ:
  ```
  TG_APP_ID="12345678"
  ```

## 📊 Что дальше?

После успешного тестирования фундамента можно переходить к:

1. **Фаза 2**: Реализация Userbot (MTProto клиент)
2. **Фаза 3**: Создание Bot Interface (команды управления)
3. **Фаза 4**: Система сбора сообщений (Ingestion)

## 🔄 Обновление кода из GitHub

```bash
git pull origin main
```

## 📝 Коммит изменений

```bash
git add .
git commit -m "Описание изменений"
git push origin main
```

## 🆘 Поддержка

При возникновении проблем:
1. Проверьте логи: `docker-compose logs -f`
2. Посмотрите STATUS.md для текущего статуса
3. Создайте issue на GitHub
