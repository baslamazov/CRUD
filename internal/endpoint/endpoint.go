package endpoint

import (
	"EffectiveMobile/internal/domain/models"
	"EffectiveMobile/internal/services"
	"io"
	"net/http"
	"strconv"
)

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Endpoint struct {
	service *services.Service
}

func New(service *services.Service) *Endpoint {
	return &Endpoint{service: service}
}
func (ep *Endpoint) GetSong(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("group")
	songName := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("release_date")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Значения по умолчанию для пагинации
	page := 1

	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	songs, _ := ep.service.GetSong(r.Context(), songName, groupName, releaseDate, limit, offset)

	responseBody, err := songs.MarshalJSON()
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
	return
}

func (ep *Endpoint) DeleteSong(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("group")
	songName := r.URL.Query().Get("song")

	success, _ := ep.service.DeleteSong(r.Context(), songName, groupName)

	responseBody := []byte("success" + strconv.FormatBool(success))
	if success != true {
		http.Error(w, "Ошибка формирования ответа", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
	return
}
func (ep *Endpoint) GetLyric(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("group")
	songName := r.URL.Query().Get("song")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Значения по умолчанию для пагинации
	page := 1

	limit := 10

	// Парсим номер страницы
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	// Парсим лимит
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	lyrics, _ := ep.service.GetLyric(r.Context(), songName, groupName, limit, offset)

	// Возвращаем ответ в формате JSON
	responseBody, err := lyrics.MarshalJSON()
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
	return
}

func (ep *Endpoint) AddSong(w http.ResponseWriter, r *http.Request) {
	var input = models.Song{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	// Десериализация JSON с помощью easyjson
	err = input.UnmarshalJSON(body)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	success, _ := ep.service.NewSong(r.Context(), input)
	if success != true {
		http.Error(w, "Ошибка добавления песни", http.StatusBadRequest)

	}

	w.WriteHeader(http.StatusCreated)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("success"))
	return
}
