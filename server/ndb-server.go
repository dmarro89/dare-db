package server

import (
	"encoding/json"
	"github.com/go-while/nodare-db/database"
	"github.com/go-while/nodare-db/logger"
	"net/http"
)

const KEY_PARAM = "key"

type WebMux interface {
	CreateMux() *http.ServeMux
	HandlerGetValByKey(w http.ResponseWriter, r *http.Request)
	HandlerSet(w http.ResponseWriter, r *http.Request)
	HandlerDel(w http.ResponseWriter, r *http.Request)
}

type XNDBServer struct {
	db     *database.XDatabase
	logger *ilog.LOG
}

func NewXNDBServer(db *database.XDatabase, logger *ilog.LOG) *XNDBServer {
	return &XNDBServer{
		db:     db,
		logger: logger,
	}
}

func (srv *XNDBServer) CreateMux() *http.ServeMux {
	mux := http.NewServeMux()
	//mux.HandleFunc("POST /loglvl/{"+KEY_PARAM+"}", srv.SetLogLvl)
	mux.HandleFunc("GET /get/{"+KEY_PARAM+"}", srv.HandlerGetValByKey)
	mux.HandleFunc("POST /set", srv.HandlerSet)
	mux.HandleFunc("GET /del/{"+KEY_PARAM+"}", srv.HandlerDel)
	return mux
}

func (srv *XNDBServer) HandlerGetValByKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.PathValue(KEY_PARAM)
	if key == "" {
		w.WriteHeader(http.StatusNotAcceptable) // 406
		return
	}

	val := srv.db.Get(key)
	if val == nil {
		w.WriteHeader(http.StatusGone) // 410
		return
	}

	// response as json with KEY:VAL ??
	/*
		response, err := json.Marshal(map[string]interface{}{key: val})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	*/

	// response as json with VAL only ?
	/*
		response, err := json.Marshal([]interface{}{val})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	*/

	// response as raw plain text with VAL only
	//w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(val.(string)))
}

func (srv *XNDBServer) HandlerSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable) // 406
		return
	}

	for key, value := range data {
		err = srv.db.Set(key, value)
		if err != nil {
			srv.logger.Warn("HandlerSet err='%v'", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)

}

func (srv *XNDBServer) HandlerDel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.PathValue(KEY_PARAM)
	if key == "" {
		w.WriteHeader(http.StatusNotAcceptable) // 406
		return
	}

	err := srv.db.Del(key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (srv *XNDBServer) SetLogLvl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	key := r.PathValue(KEY_PARAM)
	if key == "" {
		w.WriteHeader(http.StatusNotAcceptable) // 406
		return
	}
	// TODO test
	srv.logger.SetLOGLEVEL(ilog.GetLOGLEVEL(key))
	return
}
