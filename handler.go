package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post BlogPost

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if post.Title == "" || post.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = post.CreatedAt

	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		http.Error(w, "Failed to marshal tags", http.StatusInternalServerError)
		return
	}

	stmt := `INSERT INTO posts (title, content, category, tags, created_at, updated_at)
	VALUES(?, ?, ?, ?, ?, ? )`

	result, err := m.DB.Exec(stmt, post.Title, post.Content, post.Category, tagsJSON, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve post ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Post created successfully"}`))
}

func (m *PostModel) UpdatePost(w http.ResponseWriter, r *http.Request) {
	var post BlogPost

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid update index", http.StatusBadRequest)
		return
	}

	post.UpdatedAt = time.Now()

	tagsJSON, err := json.Marshal(post.Tags)
	if err != nil {
		http.Error(w, "Failed to marshal tags", http.StatusInternalServerError)
		return
	}

	stmt := `UPDATE posts SET title=?, content=?, category=?, tags=?, updated_at=? WHERE id=?`
	result, err := m.DB.Exec(stmt, post.Title, post.Content, post.Category, tagsJSON, post.UpdatedAt, id)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Post not found or no changes made", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Post updated successfully"}`))
}

func (m *PostModel) GetPost(w http.ResponseWriter, r *http.Request) {
	var post BlogPost
	var tagsJSON []byte

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	stmt := `SELECT id, title, content, category, tags, created_at, updated_at FROM posts WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)

	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &tagsJSON, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Error scanning post", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(tagsJSON, &post.Tags); err != nil {
		http.Error(w, "Failed to unmarshal tags", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		http.Error(w, "Failed to encode post", http.StatusInternalServerError)
		return
	}
}

func (m *PostModel) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	stmt := `SELECT id, title, content, category, tags, created_at, updated_at FROM posts`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []BlogPost

	for rows.Next() {
		var post BlogPost
		var tagsJSON []byte

		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Category, &tagsJSON, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			http.Error(w, "Error scanning post", http.StatusInternalServerError)
			return
		}

		if err := json.Unmarshal(tagsJSON, &post.Tags); err != nil {
			http.Error(w, "Failed to unmarshal tags", http.StatusInternalServerError)
			return
		}

		posts = append(posts, post)
	}

	if len(posts) == 0 {
		http.Error(w, "No post exists", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
		return
	}
}

func (m *PostModel) DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	stmt := `DELETE FROM posts WHERE id = ?`
	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Post deleted successfully"}`))
}
