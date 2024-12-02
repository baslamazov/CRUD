package endpoint

import (
	"EffectiveMobile/internal/domain/models"
	"EffectiveMobile/internal/services"
	"io"
	"net/http"
	"strconv"
)

type Endpoint struct {
	service *services.Service
}

func New(service *services.Service) *Endpoint {
	return &Endpoint{service: service}
}

// GetSong godoc
// @Summary Получить песню
// @Description Получить песню по имени группы и имени песни
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Имя группы"
// @Param song query string false "Имя песни"
// @Param release_date query string false "Дата релиза"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Лимит"
// @Success 200 {array} models.Song
// @Failure 500 {string} string "error"
// @Router /library/song [get]
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

// DeleteSong godoc
// @Summary Удалить песню
// @Description Удалить песню по имени группы и имени песни
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Имя группы"
// @Param song query string true "Имя песни"
// @Success 200 {string} string "success"
// @Failure 404 {string} string "error"
// @Router /library/song [delete]
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

// GetLyric godoc
// @Summary Получить текст песни
// @Description Получить текст песни по имени группы и имени песни
// @Tags lyrics
// @Accept json
// @Produce json
// @Param group query string true "Имя группы"
// @Param song query string true "Имя песни"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Лимит"
// @Success 200 {array} models.Lyric
// @Failure 500 {string} string "error"
// @Router /library/lyric [get]
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

// AddSong godoc
// @Summary Добавить песню
// @Description Добавить новую песню
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Песня"
// @Success 201 {string} string "success"
// @Failure 400 {string} string "error"
// @Router /library/song [post]
func (ep *Endpoint) AddSong(w http.ResponseWriter, r *http.Request) {
	var input = models.Song{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

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
	_, err = w.Write([]byte("success"))
	if err != nil {
		return
	}
	return
}

// UpdateSong godoc
// @Summary Обновить песню
// @Description Обновить существующую песню
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string true "Имя группы"
// @Param song query string true "Имя песни"
// @Param song body models.Song true "Песня"
// @Success 201 {string} string "success"
// @Failure 400 {string} string "error"
// @Router /library/song [put]
// TODO: Вызвать существующий метод для получения списка песен, и затем вызвать метод обновления
func (ep *Endpoint) UpdateSong(w http.ResponseWriter, r *http.Request) {
	groupName := r.URL.Query().Get("group")
	songName := r.URL.Query().Get("song")

	var input = models.Song{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	err = input.UnmarshalJSON(body)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	input.Name = songName
	input.GroupName = groupName

	success, _ := ep.service.UpdateSong(r.Context(), input)
	if success != true {
		http.Error(w, "Ошибка добавления песни", http.StatusBadRequest)

	}

	w.WriteHeader(http.StatusCreated)

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte("success"))
	if err != nil {
		return
	}
	return
}
