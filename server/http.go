package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mlinhard/ga4gh-search-go/api"
	"net/http"
	"time"
)

type HttpSearchServer struct {
	version    string
	httpServer *http.Server
	search     *SearchService
}

func NewHttpSearchServer(listenAddress string) (*HttpSearchServer, error) {
	var err error
	server := new(HttpSearchServer)
	server.search, err = NewSearchService()
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	r.HandleFunc("/tables", server.handleTables).Methods("GET")
	r.HandleFunc("/table/{table_name}/info", server.handleTableInfo).Methods("GET")
	r.HandleFunc("/table/{table_name}/data", server.handleTableData).Methods("GET")
	r.HandleFunc("/search", server.handleSearch).Methods("POST")
	r.HandleFunc("/service-info", server.handleServiceInfo).Methods("GET")

	http.Handle("/", r)
	server.httpServer = &http.Server{
		Addr:    listenAddress,
		Handler: nil,
	}
	return server, nil
}

func (server *HttpSearchServer) Service() *SearchService {
	return server.search
}

func (server *HttpSearchServer) Start() {
	go server.listenAndServe()
}

func (server *HttpSearchServer) listenAndServe() {
	err := server.httpServer.ListenAndServe()
	if err != nil {
		fmt.Printf("HTTP Server listen: %v\n", err)
	}
}

func (server *HttpSearchServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()
	err := server.httpServer.Shutdown(ctx)
	if err != nil {
		fmt.Printf("HTTP Server stop: %v", err)
	}
}

func (server *HttpSearchServer) handleTables(w http.ResponseWriter, r *http.Request) {
	tables, err := server.search.Tables()
	if err != nil {
		http.Error(w, fmt.Sprintf("Processing request: %v", err), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(tables)
}

func (server *HttpSearchServer) handleTableInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableInfo, err := server.search.TableInfo(vars["table_name"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Processing request: %v", err), http.StatusInternalServerError)
	}
	if tableInfo == nil {
		http.NotFound(w, r)
	}
	json.NewEncoder(w).Encode(tableInfo)
}

func (server *HttpSearchServer) handleTableData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tableData, err := server.search.TableData(vars["table_name"])
	if err != nil {
		http.Error(w, fmt.Sprintf("Processing request: %v", err), http.StatusInternalServerError)
	}
	if tableData == nil {
		http.NotFound(w, r)
	}
	json.NewEncoder(w).Encode(tableData)
}

func (server *HttpSearchServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	searchRequest := new(api.SearchRequest)
	err := json.NewDecoder(r.Body).Decode(searchRequest)
	if err != nil {
		http.Error(w, fmt.Sprintf("Decoding request: %v", err), http.StatusBadRequest)
	} else {
		searchResponse, err := server.search.Search(searchRequest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Processing request: %v", err), http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(searchResponse)
	}
}

func (server *HttpSearchServer) handleServiceInfo(w http.ResponseWriter, r *http.Request) {
	serviceInfo, err := server.search.ServiceInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("Processing request: %v", err), http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(serviceInfo)
}
