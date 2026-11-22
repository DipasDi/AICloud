package web

import (
	"AICloud/internal/database"
	"AICloud/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// StartWebServer запускает встроенный веб-сервер для веб-интерфейса AICloud
func StartWebServer(port string) {
	if port == "" {
		port = os.Getenv("WEB_PORT")
		if port == "" {
			port = "8080"
		}
	}

	mux := http.NewServeMux()

	// Путь к директории web
	workDir, _ := os.Getwd()
	webDir := filepath.Join(workDir, "web")

	// Статические файлы (CSS, JS, Assets)
	fileServer := http.FileServer(http.Dir(webDir))
	mux.Handle("/css/", fileServer)
	mux.Handle("/js/", fileServer)
	mux.Handle("/assets/", fileServer)

	// Главная страница веб-интерфейса
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// Проверяем, существует ли файл в директории web
			targetFile := filepath.Join(webDir, filepath.Clean(r.URL.Path))
			if _, err := os.Stat(targetFile); err == nil {
				http.ServeFile(w, r, targetFile)
				return
			}
		}
		http.ServeFile(w, r, filepath.Join(webDir, "index.html"))
	})

	// API: Получение списка записей из PostgreSQL
	mux.HandleFunc("/api/items", handleGetItems)

	// API: Поиск записей
	mux.HandleFunc("/api/search", handleSearchItems)

	// API: Статистика хранилища
	mux.HandleFunc("/api/stats", handleStats)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	fmt.Printf("🌐 [AICloud Web UI] Сервер запущен на http://localhost:%s\n", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("Ошибка веб-сервера: %v\n", err)
	}
}

// handleGetItems возвращает записи пользователя из базы данных
func handleGetItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if database.ConnPSQL == nil {
		http.Error(w, `{"error":"База данных не подключена"}`, http.StatusServiceUnavailable)
		return
	}

	targetType := r.URL.Query().Get("type")
	groupName := r.URL.Query().Get("group")

	query := `
		SELECT 
			id, user_id, chat_id, target_type, file_type, 
			tg_file_id, tg_message_id, file_name, comment, 
			group_name, created_at, updated_at
		FROM items
		WHERE 1=1
	`
	var args []interface{}
	argIdx := 1

	if targetType != "" && targetType != "all" {
		query += fmt.Sprintf(" AND target_type = $%d", argIdx)
		args = append(args, targetType)
		argIdx++
	}

	if groupName != "" && groupName != "all" {
		query += fmt.Sprintf(" AND group_name = $%d", argIdx)
		args = append(args, groupName)
		argIdx++
	}

	query += " ORDER BY created_at DESC LIMIT 100;"

	rows, err := database.ConnPSQL.Query(context.Background(), query, args...)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var it models.Item
		err := rows.Scan(
			&it.ID, &it.UserID, &it.ChatID, &it.TargetType, &it.FileType,
			&it.TgFileID, &it.TgMessageID, &it.FileName, &it.Comment,
			&it.GroupName, &it.CreatedAt, &it.UpdatedAt,
		)
		if err == nil {
			items = append(items, it)
		}
	}

	if items == nil {
		items = []models.Item{}
	}

	json.NewEncoder(w).Encode(items)
}

// handleSearchItems выполняет поиск по ключевым словам или векторам
func handleSearchItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	queryText := strings.TrimSpace(r.URL.Query().Get("q"))
	if queryText == "" {
		handleGetItems(w, r)
		return
	}

	if database.ConnPSQL == nil {
		http.Error(w, `{"error":"База данных не подключена"}`, http.StatusServiceUnavailable)
		return
	}

	searchQuery := `
		SELECT 
			id, user_id, chat_id, target_type, file_type, 
			tg_file_id, tg_message_id, file_name, comment, 
			group_name, created_at, updated_at,
			ts_rank(search_vector, plainto_tsquery('russian', $1)) AS rank
		FROM items
		WHERE search_vector @@ plainto_tsquery('russian', $1)
		ORDER BY rank DESC
		LIMIT 20;
	`

	rows, err := database.ConnPSQL.Query(context.Background(), searchQuery, queryText)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var it models.Item
		err := rows.Scan(
			&it.ID, &it.UserID, &it.ChatID, &it.TargetType, &it.FileType,
			&it.TgFileID, &it.TgMessageID, &it.FileName, &it.Comment,
			&it.GroupName, &it.CreatedAt, &it.UpdatedAt, &it.Rank,
		)
		if err == nil {
			items = append(items, it)
		}
	}

	if items == nil {
		items = []models.Item{}
	}

	json.NewEncoder(w).Encode(items)
}

// handleStats возвращает сводку по количеству сохраненных элементов
func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	stats := map[string]interface{}{
		"status":      "online",
		"db":          "PostgreSQL 17",
		"vector_dim":  768,
		"embed_model": "models/gemini-embedding-001",
		"server_time": time.Now(),
	}

	if database.ConnPSQL != nil {
		var total int64
		_ = database.ConnPSQL.QueryRow(context.Background(), "SELECT count(*) FROM items;").Scan(&total)
		stats["total_items"] = total
	}

	json.NewEncoder(w).Encode(stats)
}

// Helper: parse int64 safely
func parseInt64(val string) int64 {
	n, _ := strconv.ParseInt(val, 10, 64)
	return n
}
