package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dmarro89/dare-db/auth"
	"github.com/dmarro89/dare-db/database"
)

const KEY_PARAM = "key"
const COLLECTION_NAME_PARAM = "collectionName"

type IDare interface {
	CreateMux(auth.Authorizer, auth.Authenticator) *http.ServeMux
	HandlerGetById(w http.ResponseWriter, r *http.Request)
	HandlerSet(w http.ResponseWriter, r *http.Request)
	HandlerDelete(w http.ResponseWriter, r *http.Request)
	HandlerLogin(w http.ResponseWriter, r *http.Request)
}

type DareServer struct {
	userStore         *auth.UserStore
	collectionManager *database.CollectionManager
}

func NewDareServer(db *database.Database, userStore *auth.UserStore) *DareServer {
	collectionManager := database.NewCollectionManager()
	collectionManager.AddCollection(database.DEFAULT_COLLECTION)

	return &DareServer{
		userStore:         userStore,
		collectionManager: collectionManager,
	}
}

func (srv *DareServer) CreateMux(authorizer auth.Authorizer, authenticator auth.Authenticator) *http.ServeMux {
	mux := http.NewServeMux()

	if authorizer == nil {
		authorizer = auth.GetDefaultAuth()
	}

	if authenticator == nil {
		authenticator = auth.NewJWTAutenticatorWithUsers(srv.userStore)
	}

	middleware := auth.NewCasbinMiddleware(authorizer, authenticator)
	mux.HandleFunc(
		fmt.Sprintf(`GET /get/{%s}`, KEY_PARAM), middleware.HandleFunc(srv.HandlerGetById))
	mux.HandleFunc("POST /set", middleware.HandleFunc(srv.HandlerSet))
	mux.HandleFunc(fmt.Sprintf(`DELETE /delete/{%s}`, KEY_PARAM), middleware.HandleFunc(srv.HandlerDelete))
	mux.HandleFunc("POST /login", srv.HandlerLogin)
	mux.HandleFunc(
		fmt.Sprintf(`GET /collections/{%s}`, KEY_PARAM), middleware.HandleFunc(srv.HandlerGetCollection))
	mux.HandleFunc(
		`GET /collections`, middleware.HandleFunc(srv.HandlerGetCollections))
	mux.HandleFunc(fmt.Sprintf("POST /collections/{%s}", COLLECTION_NAME_PARAM), middleware.HandleFunc(srv.HandlerCreateCollection))
	mux.HandleFunc(fmt.Sprintf(`DELETE /collections/{%s}`, COLLECTION_NAME_PARAM), middleware.HandleFunc(srv.HandlerDeleteCollection))
	mux.HandleFunc(
		fmt.Sprintf(`GET /collections/{%s}/get/{%s}`, COLLECTION_NAME_PARAM, KEY_PARAM), middleware.HandleFunc(srv.HandlerCollectionGetById))
	mux.HandleFunc(
		fmt.Sprintf(`GET /collections/{%s}/items`, COLLECTION_NAME_PARAM), middleware.HandleFunc(srv.HandlerGetPaginatedCollectionItems))
	mux.HandleFunc(fmt.Sprintf("POST /collections/{%s}/set", COLLECTION_NAME_PARAM), middleware.HandleFunc(srv.HandlerCollectionSet))
	mux.HandleFunc(fmt.Sprintf(`DELETE /collections/{%s}/delete/{%s}`, COLLECTION_NAME_PARAM, KEY_PARAM), middleware.HandleFunc(srv.HandlerCollectionDelete))

	// Wrap the mux with the CORS handler
	corsHandler := srv.setupCORS(mux)
	// Create a new ServeMux that uses the CORS handler.
	finalMux := http.NewServeMux()
	finalMux.Handle("/", corsHandler)

	return finalMux
}

func (srv *DareServer) HandlerGetById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.PathValue(KEY_PARAM)
	if key == "" {
		http.Error(w, `url path param "key" cannot be empty`, http.StatusBadRequest)
		return
	}

	val := srv.collectionManager.GetDefaultCollection().Get(key)
	if val == "" {
		http.Error(w, fmt.Sprintf(`Key "%v" not found`, key), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(map[string]string{key: val})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (srv *DareServer) HandlerCollectionGetById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collectionName := r.PathValue(COLLECTION_NAME_PARAM)
	collection, exists := srv.collectionManager.GetCollection(collectionName)
	if !exists {
		http.Error(w, fmt.Sprintf(`Collection "%s" not found`, collectionName), http.StatusNotFound)
		return
	}

	key := r.PathValue(KEY_PARAM)
	if key == "" {
		http.Error(w, `url path param "key" cannot be empty`, http.StatusBadRequest)
		return
	}

	val := collection.Get(key)
	if val == "" {
		http.Error(w, fmt.Sprintf(`Key "%v" not found`, key), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(map[string]string{key: val})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func (srv *DareServer) HandlerGetPaginatedCollectionItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the page and pageSize from query parameters
	page := parseQueryParam(r, "page", 1)
	pageSize := parseQueryParam(r, "pageSize", 10)

	collectionName := r.PathValue(COLLECTION_NAME_PARAM)
	collection, exists := srv.collectionManager.GetCollection(collectionName)
	if !exists {
		http.Error(w, fmt.Sprintf(`Collection "%s" not found`, collectionName), http.StatusNotFound)
		return
	}

	// Retrieve paginated items
	items := collection.GetAllItems()
	paginatedItems := paginateItems(items, page, pageSize)

	response, err := json.Marshal(map[string]interface{}{
		"items":    paginatedItems,
		"page":     page,
		"pageSize": pageSize,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

type Item struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func paginateItems(items map[string]string, page, pageSize int) []Item {
	keys := make([]string, 0, len(items))
	for key := range items {
		keys = append(keys, key)
	}

	totalItems := len(keys)
	startIndex := (page - 1) * pageSize
	if startIndex >= totalItems {
		return []Item{} // Return empty array if page exceeds total items
	}

	endIndex := startIndex + pageSize
	if endIndex > totalItems {
		endIndex = totalItems
	}

	// Crea una slice di Item per contenere gli elementi paginati
	paginatedItems := make([]Item, 0, endIndex-startIndex)
	for _, key := range keys[startIndex:endIndex] {
		paginatedItems = append(paginatedItems, Item{
			Key:   key,
			Value: items[key],
		})
	}

	return paginatedItems
}

func parseQueryParam(r *http.Request, key string, defaultValue int) int {
	queryValue := r.URL.Query().Get(key)
	if queryValue == "" {
		return defaultValue
	}
	parsedValue, err := strconv.Atoi(queryValue)
	if err != nil {
		return defaultValue
	}
	return parsedValue
}

func (srv *DareServer) HandlerSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON format, the body must be in the form of {\"key\": \"value\"}", http.StatusBadRequest)
		return
	}

	for key, value := range data {
		err = srv.collectionManager.GetDefaultCollection().Set(key, value)
		if err != nil {
			http.Error(w, "Error saving data", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (srv *DareServer) HandlerCollectionSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON format, the body must be in the form of {\"key\": \"value\"}", http.StatusBadRequest)
		return
	}

	collectionName := r.PathValue(COLLECTION_NAME_PARAM)
	collection, exists := srv.collectionManager.GetCollection(collectionName)
	if !exists {
		srv.collectionManager.AddCollection(collectionName)
		collection, _ = srv.collectionManager.GetCollection(collectionName)
	}

	for key, value := range data {
		err = collection.Set(key, value)
		if err != nil {
			http.Error(w, "Error saving data", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (srv *DareServer) HandlerDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.PathValue(KEY_PARAM)
	if key == "" {
		http.Error(w, `url path param "key" cannot be empty`, http.StatusBadRequest)
		return
	}

	err := srv.collectionManager.GetDefaultCollection().Delete(key)
	if err != nil {
		http.Error(w, "Error deleting data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (srv *DareServer) HandlerCollectionDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collectionName := r.PathValue(COLLECTION_NAME_PARAM)
	collection, exists := srv.collectionManager.GetCollection(collectionName)
	if !exists {
		http.Error(w, fmt.Sprintf(`Collection "%s" not found`, collectionName), http.StatusNotFound)
		return
	}

	key := r.PathValue(KEY_PARAM)
	if key == "" {
		http.Error(w, `url path param "key" cannot be empty`, http.StatusBadRequest)
		return
	}

	err := collection.Delete(key)
	if err != nil {
		http.Error(w, "Error deleting data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (srv *DareServer) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok || !srv.userStore.ValidateCredentials(username, password) {
		http.Error(w, "Unauthorized: missing or invalid credentials", http.StatusUnauthorized)
		return
	}

	authenticator := auth.NewJWTAutenticatorWithUsers(srv.userStore)
	token, err := authenticator.GenerateToken(username)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	srv.userStore.SaveToken(username, token)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})

	if err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

func (srv *DareServer) HandlerGetCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collectionName := r.PathValue(COLLECTION_NAME_PARAM)
	collection, exists := srv.collectionManager.GetCollection(collectionName)
	if !exists {
		http.Error(w, fmt.Sprintf(`Collection "%s" not found`, collectionName), http.StatusNotFound)
		return
	}

	key := r.PathValue(KEY_PARAM)
	if key == "" {
		http.Error(w, `url path param "key" cannot be empty`, http.StatusBadRequest)
		return
	}

	val := collection.Get(key)
	if val == "" {
		http.Error(w, fmt.Sprintf(`Key "%v" not found`, key), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(map[string]string{key: val})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (srv *DareServer) HandlerGetCollections(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response, err := json.Marshal(srv.collectionManager.GetCollectionNames())
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func (srv *DareServer) HandlerCreateCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collectionName := r.PathValue(COLLECTION_NAME_PARAM)
	_, exists := srv.collectionManager.GetCollection(collectionName)
	if exists {
		http.Error(w, fmt.Sprintf(`Collection "%s" already exists`, collectionName), http.StatusBadRequest)
		return
	}

	srv.collectionManager.AddCollection(collectionName)

	w.WriteHeader(http.StatusCreated)
}

func (srv *DareServer) HandlerDeleteCollection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collectionName := r.PathValue(COLLECTION_NAME_PARAM)
	_, exists := srv.collectionManager.GetCollection(collectionName)
	if !exists {
		http.Error(w, fmt.Sprintf(`Collection "%s" not exists`, collectionName), http.StatusBadRequest)
		return
	}

	srv.collectionManager.RemoveCollection(collectionName)
	w.WriteHeader(http.StatusOK)
}

func (srv *DareServer) setupCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5002") // Or "*" for all origins (less secure)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
