# 🚀 Telegram Redesign для Figma: Page 1–7 и Responsive

Готовый плагин и скрипт для автоматического создания и редизайна страниц **Page 1 - Page 7** и **Responsive** в файле Figma в фирменном стиле **Telegram Mini Apps / Telegram Web / Telegram Dark**:
[https://www.figma.com/design/T4qokfm0BwcrPBeQ1qnXJ2/](https://www.figma.com/design/T4qokfm0BwcrPBeQ1qnXJ2/)

---

## 🎨 Что создается на холсте Figma:

1. **`Page 1 — Hero & Bot Gateway (Telegram Style)`** (1440x960):
   - Верхний Telegram WebApp бар со статусом, бот-аватаром и верификацией (`✓`).
   - Баббл сообщения в стиле Telegram с презентацией бота **AICloud**.
   - Фирменные кнопки действий: `✈️ Запустить в Telegram (@dipas_di)`, `🌐 Открыть WebApp`, `📂 GitHub`.
   - Живая телеметрия (Ping: 14ms, Goroutines: 52, pgvector HNSW: Ready).
   - Bento-плитки возможностей: Whisper Turbo (0.25с), гибридный поиск Postgres, облачное хранилище TG, Inline кнопки.

2. **`Page 2 — About & Architecture (Telegram Style)`** (1440x960):
   - Профиль разработчика в стиле Telegram-канала / профиля с аватаром, бейджем верификации и ссылками.
   - Блок философии и продакшен-метрик (Go 1.23+, p99 SLA < 15ms, pgvector HNSW 768-dim).
   - 4 принципа архитектуры чистых и легковесных микросервисов.

3. **`Page 3 — Projects & Ecosystem (Telegram Style)`** (1440x1020):
   - Карточки ключевых проектов в стиле Telegram Mini App:
     - 🧠 **AICloud Vault** (AI Telegram Bot & Knowledge Hub)
     - 🎮 **Kotichi** (High-load игровая платформа, 45k+ RPS)
     - 🗄️ **UIXDB** (Визуальный воркспейс для PostgreSQL, 5ms SLA)
     - 🌐 **NETLOS** (Гибридная платформа Minecraft + REST API)
     - 🎵 **IFinderMusic** (Web Audio плеер для музыки онлайн)

4. **`Page 4 — Skills & Telemetry Matrix (Telegram Style)`** (1440x960):
   - Список навыков в формате настроек Telegram:
     - `FOCUS 2026`: gRPC, Kafka, OpenTelemetry, Kubernetes с цветными индикаторами прогресса.
     - `CORE STACK`: Go Concurrency (95%), PostgreSQL 17 pgvector (90%), Redis (88%), Docker (85%).

5. **`Page 5 — Technical Arsenal (Telegram Style)`** (1440x960):
   - Файловое хранилище Telegram Cloud с категоризацией по рантайму, базам данных и инфраструктуре.

6. **`Page 6 — Contact & Direct Chat (Telegram Style)`** (1440x960):
   - Полноценный интерфейс диалога Telegram: карточки контактов (@dipas_di, канал @DipasDEV, email, GitHub) и поле ввода сообщения со скрепкой и кнопкой отправки.

7. **`Page 7 — Voice & AI Hub (Telegram Style)`** (1440x960):
   - Интерактивный аудиоплеер голосового сообщения Telegram со звуковой волной (waveform).
   - Транскрипция Groq Whisper Turbo за 0.25 сек.
   - Бейдж классификации намерений (Qwen 3.8 / Nemotron).
   - Telegram Inline-кнопки: `[ ⏰ Напомнить ]`, `[ 📁 В папку ]`, `[ 🗑 Удалить ]`, `[ 🙈 Скрыть ]`.
   - Пайплайн гибридного поиска (tsvector 1 мс + pgvector 768 dim).

8. **`Responsive — Telegram Mini App Mobile (390x844)`**:
   - Полный мокап Telegram WebApp на iPhone 15: статус-бар iOS, навигационная панель с кнопкой «Закрыть», фильтры категорий, карточка голосовой заметки, недавние записи и плавающая Telegram `MainButton`.

---

## ⚡ Как запустить (3 простых способа)

### Способ 1: Как локальный плагин в Figma (Рекомендуется, 1 минута)
1. Откройте Figma и перейдите в ваш файл:
   `https://www.figma.com/design/T4qokfm0BwcrPBeQ1qnXJ2/`
2. Нажмите на иконку главного меню Figma (слева вверху) ➔ **Plugins** ➔ **Development** ➔ **Import plugin from manifest...**
3. Выберите файл `manifest.json` из этой папки (`figma-telegram-redesign/manifest.json`).
4. Нажмите **Run** — плагин моментально сгенерирует все 8 экранов прямо на рабочей области!

### Способ 2: Через консоль разработчика Figma
1. В открытом файле Figma нажмите комбинацию клавиш `Ctrl + Alt + I` (на Windows) или `Cmd + Option + I` (на Mac) для открытия консоли разработчика.
2. Скопируйте весь текст из файла `code.js`.
3. Вставьте его в консоль и нажмите `Enter`.

### Способ 3: Через популярный плагин Scripter
1. Установите в Figma бесплатный плагин **Scripter**.
2. Вставьте в него код из `code.js` и нажмите кнопку **Play**.

---

## ℹ️ Информация о лимитах Figma MCP

При попытке прямой записи через протокол Figma MCP сервер вернул ограничение бесплатного тарифа Figma:
> *"You've reached the Figma MCP tool call limit on the Starter plan. Upgrade your plan for more tool calls: https://www.figma.com/files/team/1564808560661609168/all-projects?upgrade=mcp_rate_limit_paywall"*

На тарифе **Starter (бесплатный)** Figma предоставляет только 20 вызовов MCP в месяц.
При обновлении команды до тарифа Professional / Education или при назначении Dev/Full места лимит автоматически расширяется до **200 вызовов в сутки**.
Созданный нами плагин `code.js` не расходует эти лимиты и выполняется непосредственно через движок Figma!
