package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dmarro89/dare-db/auth"
	"github.com/dmarro89/dare-db/database"
)

const KEY_PARAM = "key"
const DEFAULT_USER = "user"
const DEFAULT_ROLE = "user"

type IDare interface {
	CreateMux(*auth.CasbinAuth) *http.ServeMux
	HandlerGetById(w http.ResponseWriter, r *http.Request)
	HandlerSet(w http.ResponseWriter, r *http.Request)
	HandlerDelete(w http.ResponseWriter, r *http.Request)
}

type DareServer struct {
	db *database.Database
}

func NewDareServer(db *database.Database) *DareServer {
	return &DareServer{
		db: db,
	}
}

func (srv *DareServer) getDefaultAuth() *auth.CasbinAuth {
	return auth.NewCasbinAuth("../auth/test/rbac_model.conf", "../auth/test/rbac_policy.csv", auth.Users{
		DEFAULT_USER: {Roles: []string{DEFAULT_ROLE}},
	})
}

func (srv *DareServer) CreateMux(casbinAuth *auth.CasbinAuth) *http.ServeMux {
	mux := http.NewServeMux()

	if casbinAuth == nil {
		casbinAuth = srv.getDefaultAuth()
	}

	middleware := auth.NewCasbinMiddleware(casbinAuth)
	mux.HandleFunc(
		fmt.Sprintf(`GET /get/{%s}`, KEY_PARAM), middleware.HandleFunc(srv.HandlerGetById))
	mux.HandleFunc("POST /set", middleware.HandleFunc(srv.HandlerSet))
	mux.HandleFunc(fmt.Sprintf(`DELETE /delete/{%s}`, KEY_PARAM), middleware.HandleFunc(srv.HandlerDelete))
	return mux
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

	val := srv.db.Get(key)
	if val == nil {
		http.Error(w, fmt.Sprintf(`Key "%v" not found`, key), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(map[string]interface{}{key: val})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)

}

func (srv *DareServer) HandlerSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON format, the body must be in the form of {\"key\": \"value\"}", http.StatusBadRequest)
		return
	}

	for key, value := range data {
		err = srv.db.Set(key, value)
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

	err := srv.db.Delete(key)
	if err != nil {
		http.Error(w, "Error deleting data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
