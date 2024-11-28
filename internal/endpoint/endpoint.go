package endpoint

import (
	"EffectiveMobile/internal/services"
	"fmt"
	"github.com/go-chi/chi/v5"
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

func (ep *Endpoint) Info(w http.ResponseWriter, r *http.Request) {

	song := chi.URLParam(r, "song")
	group := chi.URLParam(r, "group")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := fmt.Sprintf(`{"song":"%s","group":"%s"}`, song, group)
	//w.Write([]byte(`{"users":` + strings.Join(users, ",") + `}`))
	w.Write([]byte(response))
}
func (ep *Endpoint) GetLibrary(w http.ResponseWriter, r *http.Request) {
	// Парсим параметры запроса
	// Парсим параметры запроса из строки запроса
	groupName := r.URL.Query().Get("group")
	songName := r.URL.Query().Get("song")
	releaseDate := r.URL.Query().Get("release_date")
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

	songs, _ := ep.service.GetSong(r.Context(), songName, groupName, releaseDate, page, limit, offset)

	// Возвращаем ответ в формате JSON
	responseBody, err := songs.MarshalJSON()
	if err != nil {
		http.Error(w, "Ошибка формирования ответа", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
	return
}

//func (ep *Endpoint) AddSong(w http.ResponseWriter, r *http.Request) {
//	var input = &models.Song{}
//	// Чтение тела запроса с использованием io.ReadAll
//	body, err := io.ReadAll(r.Body)
//	if err != nil {
//		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
//		return
//	}
//
//	// Десериализация JSON с помощью easyjson
//	err = input.UnmarshalJSON(body)
//	if err != nil {
//		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
//		return
//	}
//
//	apiUrl := os.Getenv("API_URL")
//
//	// Проверка наличия слеша в конце apiUrl
//	if !strings.HasSuffix(apiUrl, "/") {
//		apiUrl += "/"
//	}
//
//	// Формирование параметров запроса
//	params := url.Values{}
//	params.Add("group", input.GroupID)
//	params.Add("song", input.Name)
//
//	// Формирование полного URL запроса
//	apiEndpoint := apiUrl + "info" + "?" + params.Encode()
//
//	// Запрос к внешнему API
//	resp, err := http.Get(apiEndpoint)
//	if err != nil {
//		http.Error(w, "Ошибка запроса к внешнему API", http.StatusBadGateway)
//		return
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		http.Error(w, "Внешний API вернул ошибку", http.StatusBadGateway)
//		return
//	}
//
//	var songDetail = &models.SongDetail{}
//
//	body, err = io.ReadAll(resp.Body)
//	if err != nil {
//		http.Error(w, "Ошибка чтения ответа от внешнего API", http.StatusInternalServerError)
//		return
//	}
//
//	err = songDetail.UnmarshalJSON(body)
//	if err != nil {
//		http.Error(w, "Ошибка обработки ответа внешнего API", http.StatusInternalServerError)
//		return
//	}
//
//	// TODO: Вызвать сервис бдля работы с бд и добавить информацию
//	ep.service.SaveSong(r.Context(), input.Name, input.GroupID, songDetail.ReleaseDate, songDetail.Text, songDetail.Link)
//	//// Начало транзакции
//	//tx, err := dbpool.Begin(context.Background())
//	//if err != nil {
//	//	logrus.Errorf("Ошибка начала транзакции: %v", err)
//	//	http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
//	//	return
//	//}
//	//defer tx.Rollback(context.Background())
//	//
//	//// Добавление группы
//	//var groupID string
//	//err = tx.QueryRow(context.Background(),
//	//	`INSERT INTO groups (name) VALUES (\$1)
//	//     ON CONFLICT (name) DO UPDATE SET name=EXCLUDED.name
//	//     RETURNING id`, input.Group).Scan(&groupID)
//	//if err != nil {
//	//	logrus.Errorf("Ошибка добавления группы: %v", err)
//	//	http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
//	//	return
//	//}
//	//
//	//// Добавление песни
//	//var songID string
//	//err = tx.QueryRow(context.Background(),
//	//	`INSERT INTO songs (group_id, song, release_date, link)
//	//     VALUES (\$1, \$2, \$3, \$4) RETURNING id`,
//	//	groupID, input.Song, songDetail.ReleaseDate,
//	//	songDetail.Link).Scan(&songID)
//	//if err != nil {
//	//	logrus.Errorf("Ошибка добавления песни: %v", err)
//	//	http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
//	//	return
//	//}
//	//
//	//// Разбивка текста песни на куплеты
//	//verses := strings.Split(songDetail.Text, "\n\n")
//	//for idx, verse := range verses {
//	//	_, err = tx.Exec(context.Background(),
//	//		`INSERT INTO lyrics (song_id, group_id,
//	//         verse_number, text) VALUES (\$1, \$2, \$3, \$4)`,
//	//		songID, groupID, idx+1, verse)
//	//	if err != nil {
//	//		logrus.Errorf("Ошибка добавления текста песни: %v", err)
//	//		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
//	//		return
//	//	}
//	//}
//	//
//	//// Фиксация транзакции
//	//err = tx.Commit(context.Background())
//	//if err != nil {
//	//	logrus.Errorf("Ошибка фиксации транзакции: %v", err)
//	//	http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
//	//	return
//	//}
//
//	w.WriteHeader(http.StatusCreated)
//
//	w.Header().Set("Content-Type", "application/json")
//	w.Write([]byte("success"))
//	return
//}
