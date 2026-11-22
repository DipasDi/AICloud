const fs = require('fs');
const path = require('path');

const pluginCode = `// Telegram Redesign for AICloud & Portfolio
// Generates Page 1 to Page 7 and Responsive in Telegram Design System Style

(async () => {
  // 1. Preload fonts
  await figma.loadFontAsync({ family: "Inter", style: "Regular" });
  await figma.loadFontAsync({ family: "Inter", style: "Medium" });
  await figma.loadFontAsync({ family: "Inter", style: "Semi Bold" });
  await figma.loadFontAsync({ family: "Inter", style: "Bold" });

  // 2. Telegram Color Palette (0..1 range)
  const c = {
    bgCanvas: { r: 14/255, g: 22/255, b: 33/255 },       // #0E1621 (Telegram Dark bg)
    bgCard: { r: 23/255, g: 33/255, b: 43/255 },         // #17212B (Telegram message/card bg)
    bgCardSubtle: { r: 32/255, g: 43/255, b: 54/255 },   // #202B36 (Telegram elevated surface)
    bgInput: { r: 36/255, g: 47/255, b: 61/255 },        // #242F3D (Telegram search/input bg)
    tgBlue: { r: 36/255, g: 129/255, b: 204/255 },       // #2481CC (Official Telegram Blue)
    tgLightBlue: { r: 42/255, g: 171/255, b: 238/255 },  // #2AABEE
    tgGreen: { r: 49/255, g: 181/255, b: 69/255 },       // #31B545 (Telegram Success/Online)
    tgAmber: { r: 255/255, g: 168/255, b: 0/255 },       // #FFA800 (Stars/Warn)
    tgViolet: { r: 135/255, g: 116/255, b: 225/255 },    // #8774E1 (Telegram AI / Premium)
    textMain: { r: 1, g: 1, b: 1 },                      // #FFFFFF
    textSec: { r: 127/255, g: 145/255, b: 164/255 },     // #7F91A4 (Telegram Muted text)
    textDim: { r: 82/255, g: 103/255, b: 120/255 },      // #526778
    borderSubtle: { r: 255/255, g: 255/255, b: 255 },   // opacity applied in paint
    borderDivider: { r: 16/255, g: 23/255, b: 31/255 },  // #10171F
  };

  // Helper: Create Text Node
  function txt(content, size = 14, style = "Regular", color = c.textMain, autoResize = "WIDTH_AND_HEIGHT") {
    const t = figma.createText();
    t.fontName = { family: "Inter", style };
    t.characters = content;
    t.fontSize = size;
    t.fills = [{ type: "SOLID", color }];
    t.textAutoResize = autoResize;
    return t;
  }

  // Helper: Create Solid Paint
  function solid(color, opacity = 1) {
    return [{ type: "SOLID", color, opacity }];
  }

  // Helper: Create Auto Layout Container
  function auto(dir = "VERTICAL", spacing = 16, padding = 24, bg = c.bgCard, radius = 16) {
    const f = figma.createFrame();
    f.layoutMode = dir;
    f.primaryAxisSizingMode = "AUTO";
    f.counterAxisSizingMode = "AUTO";
    f.itemSpacing = spacing;
    if (typeof padding === "number") {
      f.paddingTop = f.paddingBottom = f.paddingLeft = f.paddingRight = padding;
    } else if (Array.isArray(padding)) {
      f.paddingTop = padding[0];
      f.paddingRight = padding[1];
      f.paddingBottom = padding[2];
      f.paddingLeft = padding[3];
    }
    if (bg) {
      f.fills = solid(bg);
    } else {
      f.fills = [];
    }
    if (radius) f.cornerRadius = radius;
    return f;
  }

  // Helper: Create Button
  function btn(label, bg = c.tgBlue, textCol = c.textMain, radius = 12) {
    const b = auto("HORIZONTAL", 8, [12, 20, 12, 20], bg, radius);
    b.primaryAxisAlignItems = "CENTER";
    b.counterAxisAlignItems = "CENTER";
    b.appendChild(txt(label, 14, "Semi Bold", textCol));
    return b;
  }

  // Helper: Create Badge/Tag
  function badge(label, bg = c.bgInput, textCol = c.tgLightBlue) {
    const b = auto("HORIZONTAL", 4, [4, 10, 4, 10], bg, 8);
    b.appendChild(txt(label, 12, "Medium", textCol));
    return b;
  }

  // Helper: Create Telegram App Top Header (1440px)
  function createTgHeader(title, subtitle = "Telegram Mini App", verified = true) {
    const h = auto("HORIZONTAL", 16, [16, 40, 16, 40], c.bgCardSubtle, 0);
    h.resize(1440, 68);
    h.layoutSizingHorizontal = "FIXED";
    h.primaryAxisAlignItems = "SPACE_BETWEEN";
    h.counterAxisAlignItems = "CENTER";

    // Left info
    const left = auto("HORIZONTAL", 12, 0, null, 0);
    left.counterAxisAlignItems = "CENTER";

    // Avatar
    const av = figma.createEllipse();
    av.resize(40, 40);
    av.fills = solid(c.tgBlue);
    left.appendChild(av);

    const titleBox = auto("VERTICAL", 2, 0, null, 0);
    const titleRow = auto("HORIZONTAL", 6, 0, null, 0);
    titleRow.counterAxisAlignItems = "CENTER";
    titleRow.appendChild(txt(title, 16, "Semi Bold", c.textMain));
    if (verified) {
      titleRow.appendChild(txt("✓", 12, "Bold", c.tgLightBlue));
    }
    titleBox.appendChild(titleRow);
    titleBox.appendChild(txt(subtitle, 12, "Regular", c.textSec));
    left.appendChild(titleBox);

    h.appendChild(left);

    // Right action buttons
    const right = auto("HORIZONTAL", 12, 0, null, 0);
    right.counterAxisAlignItems = "CENTER";
    right.appendChild(badge("● Online", c.bgInput, c.tgGreen));
    right.appendChild(badge("Go 1.23+ // pgvector", c.bgInput, c.textSec));
    right.appendChild(btn("•••", c.bgInput, c.textMain, 8));
    h.appendChild(right);

    return h;
  }

  // Find clear space on current canvas
  let startX = 0;
  for (const child of figma.currentPage.children) {
    startX = Math.max(startX, child.x + child.width);
  }
  if (startX > 0) startX += 160;

  const createdFrames = [];

  // =========================================================================
  // FRAME 1: PAGE 1 — HERO & BOT GATEWAY (1440 x 960)
  // =========================================================================
  {
    const f1 = auto("VERTICAL", 0, 0, c.bgCanvas, 0);
    f1.name = "Page 1 — Hero & Bot Gateway (Telegram Style)";
    f1.resize(1440, 960);
    f1.layoutSizingHorizontal = "FIXED";
    f1.layoutSizingVertical = "FIXED";
    f1.x = startX;
    f1.y = 0;

    f1.appendChild(createTgHeader("AICloud 🧠 // DIPAS DEV", "Персональное облако и AI-ассистент в Telegram"));

    // Body content 2-column
    const body = auto("HORIZONTAL", 40, [50, 60, 50, 60], null, 0);
    body.resize(1440, 892);
    body.layoutSizingHorizontal = "FIXED";
    body.layoutSizingVertical = "FILL";

    // Left Col: Telegram Chat Bubble & CTA
    const leftCol = auto("VERTICAL", 24, 0, null, 0);
    leftCol.layoutSizingHorizontal = "FILL";

    // Badge
    leftCol.appendChild(badge("TELEGRAM BOT & MINI APP ECOSYSTEM", c.bgCardSubtle, c.tgLightBlue));

    // Telegram Message Bubble
    const bubble = auto("VERTICAL", 12, 20, c.bgCard, 18);
    bubble.layoutSizingHorizontal = "FILL";
    bubble.strokes = solid(c.borderSubtle, 0.08);

    const bHead = auto("HORIZONTAL", 8, 0, null, 0);
    bHead.appendChild(txt("AICloud Bot", 15, "Semi Bold", c.tgLightBlue));
    bHead.appendChild(badge("bot", c.bgCardSubtle, c.textSec));
    bubble.appendChild(bHead);

    const bText = txt(
      "Привет! Сюда можно скидывать вообще всё: мысли, заметки, ссылки, файлы, фотографии и голосовые.\\n\\n" +
      "🎙 Голосовые: распознавание через Groq Whisper Turbo за 0.25с + автоанализ намерений (поиск / заметка / задача).\\n" +
      "⚡ Двухэтапный поиск: моментальный PostgreSQL tsvector (1 мс) + смысловой pgvector (768-dim Gemini).\\n" +
      "🔒 Приватность: медиафайлы хранятся в облаке Telegram, на сервере только защищенные метаданные.",
      14, "Regular", c.textMain
    );
    bubble.appendChild(bText);

    const bFooter = auto("HORIZONTAL", 8, 0, null, 0);
    bFooter.primaryAxisAlignItems = "SPACE_BETWEEN";
    bFooter.layoutSizingHorizontal = "FILL";
    bFooter.appendChild(txt("15:26 ✓✓", 12, "Regular", c.textSec));
    bubble.appendChild(bFooter);

    leftCol.appendChild(bubble);

    // CTA Action buttons
    const actions = auto("HORIZONTAL", 16, 0, null, 0);
    actions.appendChild(btn("✈️ Запустить в Telegram (@dipas_di)", c.tgBlue, c.textMain, 14));
    actions.appendChild(btn("🌐 Открыть WebApp Дашборд", c.bgCardSubtle, c.textMain, 14));
    actions.appendChild(btn("📂 GitHub", c.bgInput, c.textSec, 14));
    leftCol.appendChild(actions);

    // Live Telemetry bar
    const tele = auto("HORIZONTAL", 20, [14, 20, 14, 20], c.bgCard, 12);
    tele.counterAxisAlignItems = "CENTER";
    tele.appendChild(txt("PING: 14ms", 12, "Medium", c.tgGreen));
    tele.appendChild(txt("•", 12, "Regular", c.textDim));
    tele.appendChild(txt("GOROUTINES: 52", 12, "Medium", c.tgLightBlue));
    tele.appendChild(txt("•", 12, "Regular", c.textDim));
    tele.appendChild(txt("POSTGRES HNSW: READY", 12, "Medium", c.textMain));
    leftCol.appendChild(tele);

    body.appendChild(leftCol);

    // Right Col: Feature Bento Tiles
    const rightCol = auto("VERTICAL", 16, 0, null, 0);
    rightCol.layoutSizingHorizontal = "FILL";

    const makeTile = (icon, title, desc, tag) => {
      const card = auto("VERTICAL", 8, 18, c.bgCard, 14);
      card.layoutSizingHorizontal = "FILL";
      card.strokes = solid(c.borderSubtle, 0.08);
      const row = auto("HORIZONTAL", 10, 0, null, 0);
      row.primaryAxisAlignItems = "SPACE_BETWEEN";
      row.layoutSizingHorizontal = "FILL";
      row.counterAxisAlignItems = "CENTER";
      row.appendChild(txt(icon + "  " + title, 15, "Semi Bold", c.textMain));
      row.appendChild(badge(tag, c.bgCardSubtle, c.tgLightBlue));
      card.appendChild(row);
      card.appendChild(txt(desc, 13, "Regular", c.textSec));
      return card;
    };

    rightCol.appendChild(makeTile("🎙", "Распознавание речи & Намерений", "Groq Whisper Turbo транскрибирует за 0.25с, а Qwen 3.8 определяет: сохранить заметку, поставить задачу или найти старое.", "AI STT"));
    rightCol.appendChild(makeTile("⚡", "Двухэтапный гибридный поиск", "Сначала быстрый полнотекстовый tsvector (1 мс), а если нет точных слов — pgvector с косинусным сходством векторов Gemini.", "HYBRID"));
    rightCol.appendChild(makeTile("🎛", "Интерактивные кнопки Telegram", "Управление заметками прямо в чате: «⏰ Напомнить», «📁 В папку», «🗑 Удалить», «🙈 Скрыть».", "INLINE UI"));
    rightCol.appendChild(makeTile("🔒", "Архитектура Telegram Cloud", "Легковесный бэкенд на Go (telego + pgx/v5). Медиафайлы не нагружают сервер, хранясь в защищенном облаке TG.", "PRIVACY"));

    body.appendChild(rightCol);
    f1.appendChild(body);
    createdFrames.push(f1);
  }

  // =========================================================================
  // FRAME 2: PAGE 2 — ABOUT & ARCHITECTURE (1440 x 960)
  // =========================================================================
  {
    const f2 = auto("VERTICAL", 0, 0, c.bgCanvas, 0);
    f2.name = "Page 2 — About & Architecture (Telegram Style)";
    f2.resize(1440, 960);
    f2.layoutSizingHorizontal = "FIXED";
    f2.layoutSizingVertical = "FIXED";
    f2.x = startX + 1560;
    f2.y = 0;

    f2.appendChild(createTgHeader("О разработчике // DIPAS", "Senior / Middle+ Backend Go Developer & Cloud Engineer"));

    const body = auto("HORIZONTAL", 40, [50, 60, 50, 60], null, 0);
    body.resize(1440, 892);
    body.layoutSizingHorizontal = "FIXED";
    body.layoutSizingVertical = "FILL";

    // Left Profile Card
    const profCard = auto("VERTICAL", 20, 30, c.bgCard, 20);
    profCard.resize(460, 780);
    profCard.layoutSizingHorizontal = "FIXED";
    profCard.layoutSizingVertical = "FIXED";
    profCard.strokes = solid(c.borderSubtle, 0.08);

    const avRow = auto("HORIZONTAL", 16, 0, null, 0);
    avRow.counterAxisAlignItems = "CENTER";
    const av = figma.createEllipse();
    av.resize(72, 72);
    av.fills = solid(c.tgBlue);
    avRow.appendChild(av);

    const nameCol = auto("VERTICAL", 4, 0, null, 0);
    const nRow = auto("HORIZONTAL", 6, 0, null, 0);
    nRow.appendChild(txt("DIPAS DEV", 20, "Bold", c.textMain));
    nRow.appendChild(txt("✓", 14, "Bold", c.tgLightBlue));
    nameCol.appendChild(nRow);
    nameCol.appendChild(txt("@dipas_di • Telegram Channel @DipasDEV", 13, "Regular", c.tgLightBlue));
    nameCol.appendChild(badge("BACKEND ARCHITECT (GO)", c.bgCardSubtle, c.tgGreen));
    avRow.appendChild(nameCol);
    profCard.appendChild(avRow);

    profCard.appendChild(txt(
      "«Единственный источник знаний — это опыт, а единственный способ научиться программировать — писать программы.»\\n\\n" +
      "Специализируюсь на создании высокопроизводительных отказоустойчивых REST API, микросервисов и ботов на Go. Строгая типизация, предсказуемость, минимальный оверхед по памяти и чистая архитектура.",
      14, "Regular", c.textSec
    ));

    // Quick Stats Grid
    const statsBox = auto("VERTICAL", 10, 16, c.bgCardSubtle, 12);
    statsBox.layoutSizingHorizontal = "FILL";
    statsBox.appendChild(txt("КЛЮЧЕВЫЕ ПОКАЗАТЕЛИ В ПРОДАКШЕНЕ", 11, "Semi Bold", c.tgLightBlue));
    
    const row1 = auto("HORIZONTAL", 16, 0, null, 0);
    row1.primaryAxisAlignItems = "SPACE_BETWEEN";
    row1.layoutSizingHorizontal = "FILL";
    row1.appendChild(txt("Основной стек:", 13, "Regular", c.textSec));
    row1.appendChild(txt("Go 1.23+ / Postgres / Redis", 13, "Semi Bold", c.textMain));
    statsBox.appendChild(row1);

    const row2 = auto("HORIZONTAL", 16, 0, null, 0);
    row2.primaryAxisAlignItems = "SPACE_BETWEEN";
    row2.layoutSizingHorizontal = "FILL";
    row2.appendChild(txt("Средний p99 SLA:", 13, "Regular", c.textSec));
    row2.appendChild(txt("< 15 ms", 13, "Semi Bold", c.tgGreen));
    statsBox.appendChild(row2);

    const row3 = auto("HORIZONTAL", 16, 0, null, 0);
    row3.primaryAxisAlignItems = "SPACE_BETWEEN";
    row3.layoutSizingHorizontal = "FILL";
    row3.appendChild(txt("Векторный поиск:", 13, "Regular", c.textSec));
    row3.appendChild(txt("pgvector HNSW (768 dim)", 13, "Semi Bold", c.tgLightBlue));
    statsBox.appendChild(row3);

    profCard.appendChild(statsBox);

    profCard.appendChild(btn("💬 Написать в Telegram", c.tgBlue, c.textMain, 12));
    body.appendChild(profCard);

    // Right Architecture Principles
    const archCol = auto("VERTICAL", 16, 0, null, 0);
    archCol.layoutSizingHorizontal = "FILL";

    const makePrinciple = (num, title, desc) => {
      const p = auto("VERTICAL", 8, 20, c.bgCard, 16);
      p.layoutSizingHorizontal = "FILL";
      p.strokes = solid(c.borderSubtle, 0.08);
      const head = auto("HORIZONTAL", 10, 0, null, 0);
      head.appendChild(badge(num, c.bgCardSubtle, c.tgLightBlue));
      head.appendChild(txt(title, 16, "Semi Bold", c.textMain));
      p.appendChild(head);
      p.appendChild(txt(desc, 14, "Regular", c.textSec));
      return p;
    };

    archCol.appendChild(makePrinciple("01", "Ориентация на результат и практическую пользу", "Любой сервис должен решать конкретную проблему пользователя, работать стабильно и не требовать постоянного ручного вмешательства."));
    archCol.appendChild(makePrinciple("02", "Чистый, типизированный и предсказуемый код", "Пишу код на Go с четкой структурой пакетов (internal/cmd/pkg), изоляцией зависимостей и понятной обработкой ошибок без паник."));
    archCol.appendChild(makePrinciple("03", "Низкий футпринт ресурсов и Zero-Alloc подход", "Оптимизация пулов соединений pgx, профилирование через pprof, минимизация выделения памяти в критических путях сервиса."));
    archCol.appendChild(makePrinciple("04", "Интеграция современных AI-инструментов", "Векторные эмбеддинги Gemini, ускоренная транскрибация Groq Whisper и классификаторы намерений на быстрых открытых моделях."));

    body.appendChild(archCol);
    f2.appendChild(body);
    createdFrames.push(f2);
  }

  // =========================================================================
  // FRAME 3: PAGE 3 — PROJECTS & ECOSYSTEM (1440 x 1020)
  // =========================================================================
  {
    const f3 = auto("VERTICAL", 0, 0, c.bgCanvas, 0);
    f3.name = "Page 3 — Projects & Ecosystem (Telegram Style)";
    f3.resize(1440, 1020);
    f3.layoutSizingHorizontal = "FIXED";
    f3.layoutSizingVertical = "FIXED";
    f3.x = startX + 3120;
    f3.y = 0;

    f3.appendChild(createTgHeader("Проекты и Экосистема", "Каталог проверенных решений в стиле Telegram Mini Apps"));

    const body = auto("VERTICAL", 24, [40, 60, 40, 60], null, 0);
    body.layoutSizingHorizontal = "FILL";

    // 2-row projects grid
    const rowTop = auto("HORIZONTAL", 24, 0, null, 0);
    rowTop.layoutSizingHorizontal = "FILL";

    const makeProjectCard = (icon, type, title, desc, metricVal, metricKey, tags, ctaLabel) => {
      const card = auto("VERTICAL", 14, 20, c.bgCard, 16);
      card.layoutSizingHorizontal = "FILL";
      card.strokes = solid(c.borderSubtle, 0.08);

      const top = auto("HORIZONTAL", 8, 0, null, 0);
      top.primaryAxisAlignItems = "SPACE_BETWEEN";
      top.layoutSizingHorizontal = "FILL";
      top.appendChild(badge(type, c.bgCardSubtle, c.tgLightBlue));
      top.appendChild(badge("ACTIVE", c.bgCardSubtle, c.tgGreen));
      card.appendChild(top);

      const tRow = auto("HORIZONTAL", 8, 0, null, 0);
      tRow.appendChild(txt(icon + " " + title, 18, "Bold", c.textMain));
      card.appendChild(tRow);

      card.appendChild(txt(desc, 13, "Regular", c.textSec));

      // Metrics
      const mBox = auto("HORIZONTAL", 16, [10, 14, 10, 14], c.bgCardSubtle, 10);
      mBox.layoutSizingHorizontal = "FILL";
      mBox.primaryAxisAlignItems = "SPACE_BETWEEN";
      const mLeft = auto("VERTICAL", 2, 0, null, 0);
      mLeft.appendChild(txt(metricVal, 16, "Bold", c.tgLightBlue));
      mLeft.appendChild(txt(metricKey, 11, "Regular", c.textSec));
      mBox.appendChild(mLeft);
      card.appendChild(mBox);

      // Tags row
      const tagRow = auto("HORIZONTAL", 6, 0, null, 0);
      tags.forEach(t => tagRow.appendChild(badge(t, c.bgInput, c.textMain)));
      card.appendChild(tagRow);

      // Action button
      card.appendChild(btn(ctaLabel, c.tgBlue, c.textMain, 10));

      return card;
    };

    rowTop.appendChild(makeProjectCard("🧠", "AI TELEGRAM BOT", "AICloud Vault", "Персональное облако и база знаний. Понимает смысл голосовых заметок через Groq Whisper Turbo и находит всё в pgvector.", "1 ms / 768-dim", "Гибридный поиск Postgres", ["Go", "telego", "pgvector", "Whisper", "Gemini"], "🚀 Запустить в Telegram"));
    rowTop.appendChild(makeProjectCard("🎮", "HIGH-LOAD SOCIAL", "Kotichi", "Мобильная социальная сеть и игровой хаб геймеров. Единый профиль, поиск Duo, мессенджер реального времени.", "45k+ RPS", "Throughput SLA 99.98%", ["Go", "Gin", "Redis", "Docker", "Goroutines"], "🔗 Подробнее о проекте"));
    rowTop.appendChild(makeProjectCard("🗄️", "DATABASE VISUALIZER", "UIXDB", "Визуальный воркспейс для PostgreSQL: ERD-схемы, транзакционный DDL, SQL-консоль без тяжелых зависимостей.", "5 ms SLA", "0 Frontend Dependencies", ["Go", "PostgreSQL", "Vanilla JS"], "📂 Открыть на GitHub"));
    body.appendChild(rowTop);

    const rowBottom = auto("HORIZONTAL", 24, 0, null, 0);
    rowBottom.layoutSizingHorizontal = "FILL";
    rowBottom.appendChild(makeProjectCard("🌐", "HYBRID PLATFORM", "NETLOS", "Гибридный сервер Minecraft и веб-экосистема с синхронизацией данных и микросервисной архитектурой в Docker.", "Real-time", "RCON & MariaDB Sync", ["Go", "Docker", "REST API", "Velocity"], "🔗 Обзор архитектуры"));
    rowBottom.appendChild(makeProjectCard("🎵", "WEB AUDIO STREAM", "IFinderMusic", "Современный веб-плеер для стриминга музыки онлайн. Мгновенный поиск, отзывчивый эквалайзер и мобильная адаптивность.", "music.dipass.ru", "Real-time Web Audio API", ["JavaScript", "Web Audio", "Streaming"], "🎧 Слушать музыку онлайн"));
    body.appendChild(rowBottom);

    f3.appendChild(body);
    createdFrames.push(f3);
  }

  // =========================================================================
  // FRAME 4: PAGE 4 — SKILLS & TELEMETRY MATRIX (1440 x 960)
  // =========================================================================
  {
    const f4 = auto("VERTICAL", 0, 0, c.bgCanvas, 0);
    f4.name = "Page 4 — Skills & Telemetry Matrix (Telegram Style)";
    f4.resize(1440, 960);
    f4.layoutSizingHorizontal = "FIXED";
    f4.layoutSizingVertical = "FIXED";
    f4.x = startX + 4680;
    f4.y = 0;

    f4.appendChild(createTgHeader("Телеметрия стека и Прогресс", "Системный радар компетенций // Go 1.23 Linux Cluster"));

    const body = auto("HORIZONTAL", 32, [40, 60, 40, 60], null, 0);
    body.resize(1440, 892);
    body.layoutSizingHorizontal = "FIXED";
    body.layoutSizingVertical = "FILL";

    const makeSkillSection = (title, badgeText, items, badgeColor) => {
      const col = auto("VERTICAL", 14, 24, c.bgCard, 18);
      col.layoutSizingHorizontal = "FILL";
      col.strokes = solid(c.borderSubtle, 0.08);

      const head = auto("HORIZONTAL", 10, 0, null, 0);
      head.primaryAxisAlignItems = "SPACE_BETWEEN";
      head.layoutSizingHorizontal = "FILL";
      head.appendChild(txt(title, 16, "Bold", c.textMain));
      head.appendChild(badge(badgeText, c.bgCardSubtle, badgeColor));
      col.appendChild(head);

      items.forEach(it => {
        const row = auto("VERTICAL", 6, 12, c.bgCardSubtle, 12);
        row.layoutSizingHorizontal = "FILL";

        const topRow = auto("HORIZONTAL", 8, 0, null, 0);
        topRow.primaryAxisAlignItems = "SPACE_BETWEEN";
        topRow.layoutSizingHorizontal = "FILL";
        topRow.appendChild(txt(it.name, 14, "Semi Bold", c.textMain));
        topRow.appendChild(txt(it.percent + "%", 14, "Bold", badgeColor));
        row.appendChild(topRow);

        // Progress bar container
        const bar = auto("HORIZONTAL", 0, 0, c.bgInput, 4);
        bar.resize(100, 6);
        bar.layoutSizingHorizontal = "FILL";
        // Fill portion
        const fill = figma.createFrame();
        fill.resize(Math.round(400 * (it.percent / 100)), 6);
        fill.fills = solid(badgeColor);
        fill.cornerRadius = 4;
        bar.appendChild(fill);
        row.appendChild(bar);

        row.appendChild(txt(it.desc, 12, "Regular", c.textSec));
        col.appendChild(row);
      });

      return col;
    };

    body.appendChild(makeSkillSection("В АКТИВНОМ ИЗУЧЕНИИ", "FOCUS 2026", [
      { name: "gRPC & Protocol Buffers", percent: 72, desc: "Двунаправленный streaming, Protobuf v3 схемы, grpc-gateway, интерцепторы." },
      { name: "Apache Kafka & Event-Driven", percent: 65, desc: "Топики, партиции, Consumer Groups, Exactly-Once семантика, DLQ, шина событий." },
      { name: "Prometheus & OpenTelemetry", percent: 60, desc: "Сквозной distributed tracing, сбор RED-метрик в Go, замеры p95/p99." },
      { name: "Kubernetes & Helm", percent: 48, desc: "Pods, Deployments, Services, Helm-чарты, Ingress, автомасштабирование HPA." }
    ], c.tgLightBlue));

    body.appendChild(makeSkillSection("ОСВОЕННАЯ БАЗА (В ПРОДЕ)", "CORE STACK", [
      { name: "Go (Golang Core & Concurrency)", percent: 95, desc: "Горутины, каналы, sync/atomic, context, Gin/fasthttp, профилирование pprof." },
      { name: "PostgreSQL 17 & pgvector", percent: 90, desc: "B-Tree/GIN индексы, HNSW векторный поиск, транзакции, pgx connection pool." },
      { name: "Redis In-Memory Архитектура", percent: 88, desc: "TTL стратегии кэширования, Redlock распределенные блокировки, Pub/Sub." },
      { name: "Docker & Docker Compose", percent: 85, desc: "Multi-stage сборки alpine/scratch (~15MB), изоляция микросервисов, сети." }
    ], c.tgGreen));

    f4.appendChild(body);
    createdFrames.push(f4);
  }

  // =========================================================================
  // FRAME 5: PAGE 5 — TECHNICAL ARSENAL (1440 x 960)
  // =========================================================================
  {
    const f5 = auto("VERTICAL", 0, 0, c.bgCanvas, 0);
    f5.name = "Page 5 — Technical Arsenal (Telegram Style)";
    f5.resize(1440, 960);
    f5.layoutSizingHorizontal = "FIXED";
    f5.layoutSizingVertical = "FIXED";
    f5.x = startX + 6240;
    f5.y = 0;

    f5.appendChild(createTgHeader("Технический Арсенал // Telegram Vault", "Структурированная база знаний и инструментов"));

    const body = auto("HORIZONTAL", 32, [40, 60, 40, 60], null, 0);
    body.resize(1440, 892);
    body.layoutSizingHorizontal = "FIXED";
    body.layoutSizingVertical = "FILL";

    const makeArsenalCol = (icon, category, items, tag) => {
      const col = auto("VERTICAL", 14, 24, c.bgCard, 18);
      col.layoutSizingHorizontal = "FILL";
      col.strokes = solid(c.borderSubtle, 0.08);

      const top = auto("HORIZONTAL", 8, 0, null, 0);
      top.primaryAxisAlignItems = "SPACE_BETWEEN";
      top.layoutSizingHorizontal = "FILL";
      top.counterAxisAlignItems = "CENTER";
      top.appendChild(txt(icon + "  " + category, 16, "Bold", c.textMain));
      top.appendChild(badge(tag, c.bgCardSubtle, c.tgLightBlue));
      col.appendChild(top);

      items.forEach(it => {
        const itemBox = auto("HORIZONTAL", 12, 14, c.bgCardSubtle, 12);
        itemBox.layoutSizingHorizontal = "FILL";
        itemBox.counterAxisAlignItems = "CENTER";

        const dot = figma.createEllipse();
        dot.resize(8, 8);
        dot.fills = solid(c.tgLightBlue);
        itemBox.appendChild(dot);

        const txtCol = auto("VERTICAL", 2, 0, null, 0);
        txtCol.layoutSizingHorizontal = "FILL";
        txtCol.appendChild(txt(it.name, 14, "Semi Bold", c.textMain));
        txtCol.appendChild(txt(it.detail, 12, "Regular", c.textSec));
        itemBox.appendChild(txtCol);

        col.appendChild(itemBox);
      });

      return col;
    };

    body.appendChild(makeArsenalCol("⚡", "GOLANG CORE", [
      { name: "Goroutines & Channels", detail: "Конкурентные пайплайны, worker pools, select timeouts" },
      { name: "Context Lifecycle", detail: "WithCancel, WithTimeout, deadline propagation" },
      { name: "High-Performance Routers", detail: "Gin, fasthttp, Gorilla Mux, net/http" },
      { name: "Memory & Allocations", detail: "pprof, бенчмаркинг, sync.Pool, zero-alloc" },
      { name: "Testing & Mocks", detail: "Unit-тесты, testify, testcontainers для интеграции" }
    ], "RUNTIME"));

    body.appendChild(makeArsenalCol("🗄️", "DATABASES & PGVECTOR", [
      { name: "PostgreSQL 17 & pgvector", detail: "Векторные 768-dim эмбеддинги, HNSW быстрый индекс" },
      { name: "Connection Pooling", detail: "pgx/v5 connection pool, эффективный мультиплексинг" },
      { name: "Composite Indexes", detail: "B-Tree, GIN, BRIN, анализ планов через EXPLAIN ANALYZE" },
      { name: "Redis In-Memory", detail: "TTL стратегии, Redlock блокировки, Pub/Sub шина" },
      { name: "DB Migrations", detail: "golang-migrate, goose, версионирование схемы" }
    ], "STORAGE"));

    body.appendChild(makeArsenalCol("🛠️", "INFRASTRUCTURE & DEVOPS", [
      { name: "Docker Multi-Stage", detail: "Миниатюрные образы scratch/alpine (~15MB per binary)" },
      { name: "Docker Compose", detail: "Локальный микросервисный кластер и тестовые БД" },
      { name: "Linux Server Operations", detail: "systemd юниты, bash-скрипты, сетевая диагностика" },
      { name: "Git & GitHub Actions", detail: "Автоматический CI/CD, линтинг, тестирование PR" }
    ], "DEVOPS"));

    f5.appendChild(body);
    createdFrames.push(f5);
  }

  // =========================================================================
  // FRAME 6: PAGE 6 — CONTACT & DIRECT CHAT (1440 x 960)
  // =========================================================================
  {
    const f6 = auto("VERTICAL", 0, 0, c.bgCanvas, 0);
    f6.name = "Page 6 — Contact & Direct Chat (Telegram Style)";
    f6.resize(1440, 960);
    f6.layoutSizingHorizontal = "FIXED";
    f6.layoutSizingVertical = "FIXED";
    f6.x = startX + 7800;
    f6.y = 0;

    f6.appendChild(createTgHeader("Связаться // Прямой диалог", "Telegram Chat Room Interface"));

    const body = auto("HORIZONTAL", 40, [40, 60, 40, 60], null, 0);
    body.resize(1440, 892);
    body.layoutSizingHorizontal = "FIXED";
    body.layoutSizingVertical = "FILL";

    // Left: Chat Window (Telegram Mockup)
    const chatWin = auto("VERTICAL", 16, 24, c.bgCard, 20);
    chatWin.resize(760, 780);
    chatWin.layoutSizingHorizontal = "FIXED";
    chatWin.layoutSizingVertical = "FIXED";
    chatWin.strokes = solid(c.borderSubtle, 0.08);

    // Chat Header
    const chHead = auto("HORIZONTAL", 12, [0, 0, 16, 0], null, 0);
    chHead.layoutSizingHorizontal = "FILL";
    chHead.counterAxisAlignItems = "CENTER";
    const chAv = figma.createEllipse();
    chAv.resize(44, 44);
    chAv.fills = solid(c.tgBlue);
    chHead.appendChild(chAv);

    const chTitle = auto("VERTICAL", 2, 0, null, 0);
    chTitle.appendChild(txt("DIPAS (Николай)", 16, "Bold", c.textMain));
    chTitle.appendChild(txt("был(а) только что • @dipas_di", 12, "Regular", c.tgGreen));
    chHead.appendChild(chTitle);
    chatWin.appendChild(chHead);

    // Message 1 from user
    const msg1 = auto("VERTICAL", 6, 14, c.bgCardSubtle, 16);
    msg1.appendChild(txt("Здравствуйте! Ищем опытного backend-разработчика на Go для высоконагруженного сервиса и интеграции с Telegram.", 14, "Regular", c.textMain));
    msg1.appendChild(txt("15:30 ✓✓", 11, "Regular", c.textSec));
    chatWin.appendChild(msg1);

    // Reply from Dipas
    const msg2 = auto("VERTICAL", 6, 14, c.bgInput, 16);
    msg2.appendChild(txt("Привет! Готов обсудить архитектуру задачи, проектирование API и микросервисов. Напишите мне в личные сообщения или оставьте заявку справа — оперативно выйду на связь!", 14, "Regular", c.textMain));
    msg2.appendChild(txt("15:31 ✓✓", 11, "Regular", c.textSec));
    chatWin.appendChild(msg2);

    // Input Bar at bottom of chat
    const inputBar = auto("HORIZONTAL", 12, [12, 16, 12, 16], c.bgCardSubtle, 24);
    inputBar.layoutSizingHorizontal = "FILL";
    inputBar.counterAxisAlignItems = "CENTER";
    inputBar.appendChild(txt("📎", 16, "Regular", c.textSec));
    const inTxt = txt("Напишите сообщение...", 14, "Regular", c.textSec);
    inTxt.layoutSizingHorizontal = "FILL";
    inputBar.appendChild(inTxt);
    inputBar.appendChild(txt("🎙", 16, "Regular", c.textSec));
    inputBar.appendChild(btn("✈️", c.tgBlue, c.textMain, 20));
    chatWin.appendChild(inputBar);

    body.appendChild(chatWin);

    // Right: Direct Channels
    const rightCards = auto("VERTICAL", 20, 0, null, 0);
    rightCards.layoutSizingHorizontal = "FILL";

    const makeContactCard = (icon, label, value, sub, linkBtn) => {
      const cBox = auto("VERTICAL", 10, 20, c.bgCard, 16);
      cBox.layoutSizingHorizontal = "FILL";
      cBox.strokes = solid(c.borderSubtle, 0.08);

      const r = auto("HORIZONTAL", 10, 0, null, 0);
      r.counterAxisAlignItems = "CENTER";
      r.appendChild(txt(icon, 20, "Regular", c.textMain));
      const col = auto("VERTICAL", 2, 0, null, 0);
      col.appendChild(txt(label, 12, "Medium", c.textSec));
      col.appendChild(txt(value, 16, "Semi Bold", c.textMain));
      r.appendChild(col);
      cBox.appendChild(r);

      cBox.appendChild(txt(sub, 12, "Regular", c.textSec));
      cBox.appendChild(btn(linkBtn, c.bgCardSubtle, c.tgLightBlue, 10));
      return cBox;
    };

    rightCards.appendChild(makeContactCard("✈️", "TELEGRAM", "@dipas_di", "Прямой чат для быстрого ответа и обсуждения", "Открыть чат @dipas_di"));
    rightCards.appendChild(makeContactCard("📢", "TELEGRAM КАНАЛ", "@DipasDEV", "Мысли об архитектуре, Go, базы данных и AI", "Подписаться на канал"));
    rightCards.appendChild(makeContactCard("✉️", "EMAIL", "dipasdarp@gmail.com", "Для официальных запросов и предложений о работе", "Скопировать email"));
    rightCards.appendChild(makeContactCard("🐙", "GITHUB", "github.com/DipasDi", "Исходный код открытых проектов и pet-репозиториев", "Перейти на GitHub"));

    body.appendChild(rightCards);
    f6.appendChild(body);
    createdFrames.push(f6);
  }

  // =========================================================================
  // FRAME 7: PAGE 7 — VOICE & AI HUB (1440 x 960)
  // =========================================================================
  {
    const f7 = auto("VERTICAL", 0, 0, c.bgCanvas, 0);
    f7.name = "Page 7 — Voice & AI Hub (Telegram Style)";
    f7.resize(1440, 960);
    f7.layoutSizingHorizontal = "FIXED";
    f7.layoutSizingVertical = "FIXED";
    f7.x = startX + 9360;
    f7.y = 0;

    f7.appendChild(createTgHeader("AI & Голосовой Хаб // Архитектура Бот-ядра", "Groq Whisper Turbo + Qwen Intent Classification + pgvector"));

    const body = auto("HORIZONTAL", 36, [40, 60, 40, 60], null, 0);
    body.resize(1440, 892);
    body.layoutSizingHorizontal = "FIXED";
    body.layoutSizingVertical = "FILL";

    // Left: Voice Player Mockup
    const voiceCard = auto("VERTICAL", 20, 28, c.bgCard, 20);
    voiceCard.resize(680, 780);
    voiceCard.layoutSizingHorizontal = "FIXED";
    voiceCard.layoutSizingVertical = "FIXED";
    voiceCard.strokes = solid(c.borderSubtle, 0.08);

    voiceCard.appendChild(badge("ИНТЕРАКТИВНЫЙ ГОЛОСОВОЙ ПЛЕЕР В ЧАТЕ", c.bgCardSubtle, c.tgLightBlue));

    // Waveform player container
    const waveBox = auto("HORIZONTAL", 16, 20, c.bgCardSubtle, 16);
    waveBox.layoutSizingHorizontal = "FILL";
    waveBox.counterAxisAlignItems = "CENTER";

    // Play button circle
    const playCircle = figma.createEllipse();
    playCircle.resize(48, 48);
    playCircle.fills = solid(c.tgBlue);
    waveBox.appendChild(playCircle);

    // Audio Bars
    const barsContainer = auto("HORIZONTAL", 4, 0, null, 0);
    barsContainer.counterAxisAlignItems = "CENTER";
    const heights = [12, 24, 18, 32, 28, 40, 20, 14, 36, 26, 16, 30, 22, 34, 18, 28, 38, 20, 12, 24];
    heights.forEach(h => {
      const b = figma.createRectangle();
      b.resize(3, h);
      b.fills = solid(c.tgLightBlue);
      b.cornerRadius = 2;
      barsContainer.appendChild(b);
    });
    waveBox.appendChild(barsContainer);

    waveBox.appendChild(txt("0:18", 14, "Medium", c.textSec));
    voiceCard.appendChild(waveBox);

    // Transcription text
    voiceCard.appendChild(txt("РАСПОЗНАННЫЙ ТЕКСТ (GROQ WHISPER TURBO — 0.25S):", 12, "Semi Bold", c.textSec));
    const transBox = auto("VERTICAL", 8, 16, c.bgInput, 12);
    transBox.layoutSizingHorizontal = "FILL";
    transBox.appendChild(txt("«Запиши код от домофона 4210 и напомни завтра в 10 утра позвонить в клинику»", 15, "Medium", c.textMain));
    voiceCard.appendChild(transBox);

    // Classification Result
    voiceCard.appendChild(txt("АНАЛИЗ НАМЕРЕНИЙ И МАРШРУТИЗАЦИЯ:", 12, "Semi Bold", c.textSec));
    const intentBox = auto("VERTICAL", 10, 16, c.bgCardSubtle, 12);
    intentBox.layoutSizingHorizontal = "FILL";
    intentBox.appendChild(badge("🎯 НАМЕРЕНИЕ: ЗАМЕТКА + НАПОМИНАНИЕ (Уверенность 99.4%)", c.bgInput, c.tgGreen));
    intentBox.appendChild(txt("Бот распознал дату («завтра 10:00») и автоматически запланировал уведомление, сохранив текст в базу.", 13, "Regular", c.textSec));
    voiceCard.appendChild(intentBox);

    // Telegram Inline Buttons
    voiceCard.appendChild(txt("ИНТЕРАКТИВНЫЕ КНОПКИ В ЧАТЕ:", 12, "Semi Bold", c.textSec));
    const btnRow = auto("HORIZONTAL", 10, 0, null, 0);
    btnRow.appendChild(btn("⏰ Напомнить", c.bgInput, c.textMain, 10));
    btnRow.appendChild(btn("📁 В папку", c.bgInput, c.textMain, 10));
    btnRow.appendChild(btn("🗑 Удалить", c.bgInput, c.textSec, 10));
    btnRow.appendChild(btn("🙈 Скрыть", c.bgInput, c.textSec, 10));
    voiceCard.appendChild(btnRow);

    body.appendChild(voiceCard);

    // Right: Vector & Database Pipeline
    const pipeCol = auto("VERTICAL", 16, 0, null, 0);
    pipeCol.layoutSizingHorizontal = "FILL";

    const makePipeStep = (step, title, detail, tech) => {
      const p = auto("VERTICAL", 8, 20, c.bgCard, 16);
      p.layoutSizingHorizontal = "FILL";
      p.strokes = solid(c.borderSubtle, 0.08);

      const r = auto("HORIZONTAL", 10, 0, null, 0);
      r.primaryAxisAlignItems = "SPACE_BETWEEN";
      r.layoutSizingHorizontal = "FILL";
      r.appendChild(txt(step + "  " + title, 15, "Semi Bold", c.textMain));
      r.appendChild(badge(tech, c.bgCardSubtle, c.tgLightBlue));
      p.appendChild(r);

      p.appendChild(txt(detail, 13, "Regular", c.textSec));
      return p;
    };

    pipeCol.appendChild(makePipeStep("1.", "Полнотекстовый поиск Postgres", "Быстрый поиск по лексемам через tsvector и GIN-индекс. Скорость 1 мс, нулевой расход AI токенов.", "Postgres tsvector"));
    pipeCol.appendChild(makePipeStep("2.", "Смысловая векторизация", "Если точных совпадений нет, запрос кодируется в 768-мерный вектор через Google Gemini Embedding.", "Gemini 768-dim"));
    pipeCol.appendChild(makePipeStep("3.", "Векторный поиск HNSW", "Индекс pgvector находит ближайшие по смыслу заметки по косинусному расстоянию даже при опечатках.", "pgvector HNSW"));
    pipeCol.appendChild(makePipeStep("4.", "Прямой переход к сообщению", "Если найдена конкретная запись, бот делает reply на оригинальное сообщение в Telegram для моментального перехода.", "Telegram Deep Link"));

    body.appendChild(pipeCol);
    f7.appendChild(body);
    createdFrames.push(f7);
  }

  // =========================================================================
  // FRAME 8: RESPONSIVE — TELEGRAM MINI APP MOBILE (390 x 844)
  // =========================================================================
  {
    const f8 = auto("VERTICAL", 0, 0, c.bgCanvas, 32);
    f8.name = "Responsive — Telegram Mini App Mobile (390x844)";
    f8.resize(390, 844);
    f8.layoutSizingHorizontal = "FIXED";
    f8.layoutSizingVertical = "FIXED";
    f8.x = startX + 10920;
    f8.y = 0;
    f8.strokes = solid(c.borderSubtle, 0.12);

    // iOS Status Bar Mockup
    const statusBar = auto("HORIZONTAL", 0, [12, 20, 8, 20], null, 0);
    statusBar.resize(390, 44);
    statusBar.layoutSizingHorizontal = "FIXED";
    statusBar.primaryAxisAlignItems = "SPACE_BETWEEN";
    statusBar.counterAxisAlignItems = "CENTER";
    statusBar.appendChild(txt("9:41", 14, "Semi Bold", c.textMain));
    statusBar.appendChild(txt("5G  100%", 12, "Medium", c.textMain));
    f8.appendChild(statusBar);

    // Telegram Mini App Navigation Bar
    const navBar = auto("HORIZONTAL", 0, [8, 16, 12, 16], c.bgCardSubtle, 0);
    navBar.resize(390, 52);
    navBar.layoutSizingHorizontal = "FIXED";
    navBar.primaryAxisAlignItems = "SPACE_BETWEEN";
    navBar.counterAxisAlignItems = "CENTER";
    navBar.appendChild(txt("Закрыть", 15, "Regular", c.tgLightBlue));

    const navTitle = auto("VERTICAL", 2, 0, null, 0);
    navTitle.counterAxisAlignItems = "CENTER";
    navTitle.appendChild(txt("AICloud Vault", 15, "Semi Bold", c.textMain));
    navTitle.appendChild(txt("bot by @dipas_di", 11, "Regular", c.textSec));
    navBar.appendChild(navTitle);

    navBar.appendChild(txt("•••", 16, "Bold", c.tgLightBlue));
    f8.appendChild(navBar);

    // Search Bar
    const searchBar = auto("HORIZONTAL", 8, [8, 12, 8, 12], c.bgInput, 10);
    searchBar.resize(358, 38);
    searchBar.counterAxisAlignItems = "CENTER";
    searchBar.x = 16;
    searchBar.appendChild(txt("🔍", 13, "Regular", c.textSec));
    searchBar.appendChild(txt("Поиск по смыслу или тексту...", 13, "Regular", c.textSec));
    
    const sbWrap = auto("VERTICAL", 0, [12, 16, 8, 16], null, 0);
    sbWrap.layoutSizingHorizontal = "FILL";
    sbWrap.appendChild(searchBar);
    f8.appendChild(sbWrap);

    // Category Filter Chips
    const chipsRow = auto("HORIZONTAL", 8, [0, 16, 8, 16], null, 0);
    chipsRow.layoutSizingHorizontal = "FILL";
    chipsRow.appendChild(badge("Все 34", c.tgBlue, c.textMain));
    chipsRow.appendChild(badge("📝 Заметки 18", c.bgCardSubtle, c.textSec));
    chipsRow.appendChild(badge("📁 Файлы 6", c.bgCardSubtle, c.textSec));
    chipsRow.appendChild(badge("🎙 Голос 7", c.bgCardSubtle, c.textSec));
    f8.appendChild(chipsRow);

    // Active Bento Voice Card (Mobile)
    const mVoiceWrap = auto("VERTICAL", 0, [4, 16, 8, 16], null, 0);
    mVoiceWrap.layoutSizingHorizontal = "FILL";

    const mVoiceCard = auto("VERTICAL", 10, 14, c.bgCard, 14);
    mVoiceCard.layoutSizingHorizontal = "FILL";
    mVoiceCard.strokes = solid(c.borderSubtle, 0.08);

    const mVoiceTop = auto("HORIZONTAL", 8, 0, null, 0);
    mVoiceTop.primaryAxisAlignItems = "SPACE_BETWEEN";
    mVoiceTop.layoutSizingHorizontal = "FILL";
    mVoiceTop.appendChild(txt("🎙 Голосовая заметка", 13, "Semi Bold", c.textMain));
    mVoiceTop.appendChild(badge("Whisper 0.2s", c.bgCardSubtle, c.tgGreen));
    mVoiceCard.appendChild(mVoiceTop);

    mVoiceCard.appendChild(txt("«Запиши код от домофона 4210»", 13, "Regular", c.textSec));
    mVoiceCard.appendChild(badge("🎯 Намерение: Заметка", c.bgInput, c.tgLightBlue));

    mVoiceWrap.appendChild(mVoiceCard);
    f8.appendChild(mVoiceWrap);

    // Recent Notes List
    const mListWrap = auto("VERTICAL", 8, [4, 16, 12, 16], null, 0);
    mListWrap.layoutSizingHorizontal = "FILL";
    mListWrap.appendChild(txt("НЕДАВНИЕ ЗАПИСИ", 12, "Semi Bold", c.textSec));

    const makeMobileItem = (icon, title, desc, date) => {
      const it = auto("HORIZONTAL", 12, 12, c.bgCard, 12);
      it.layoutSizingHorizontal = "FILL";
      it.counterAxisAlignItems = "CENTER";
      it.appendChild(txt(icon, 18, "Regular", c.textMain));
      const col = auto("VERTICAL", 2, 0, null, 0);
      col.layoutSizingHorizontal = "FILL";
      col.appendChild(txt(title, 13, "Semi Bold", c.textMain));
      col.appendChild(txt(desc, 11, "Regular", c.textSec));
      it.appendChild(col);
      it.appendChild(txt(date, 11, "Regular", c.textDim));
      return it;
    };

    mListWrap.appendChild(makeMobileItem("📝", "План спринта Go", "gRPC streaming + Kafka Event Bus", "14:20"));
    mListWrap.appendChild(makeMobileItem("📁", "architecture_schema.sql", "PostgreSQL DDL + pgvector index", "Вчера"));
    mListWrap.appendChild(makeMobileItem("🔗", "music.dipass.ru", "Онлайн плеер IFinderMusic", "20 ноя"));
    f8.appendChild(mListWrap);

    // Telegram MainButton at Bottom
    const bottomBar = auto("VERTICAL", 8, [8, 16, 16, 16], c.bgCardSubtle, 0);
    bottomBar.resize(390, 78);
    bottomBar.layoutSizingHorizontal = "FIXED";
    const mainBtn = btn("➕ ДОБАВИТЬ ЗАМЕТКУ / ГОЛОС", c.tgBlue, c.textMain, 12);
    mainBtn.layoutSizingHorizontal = "FILL";
    bottomBar.appendChild(mainBtn);
    f8.appendChild(bottomBar);

    createdFrames.push(f8);
  }

  // Zoom to fit newly created frames
  figma.viewport.scrollAndZoomIntoView(createdFrames);

  figma.notify("✅ Дизайн Page 1-7 и Responsive в стиле Telegram успешно создан!");
  return { success: true, count: createdFrames.length, frameIds: createdFrames.map(f => f.id) };
})();
`;

fs.writeFileSync(path.join(__dirname, 'code.js'), pluginCode, 'utf8');
console.log('Successfully generated code.js');
