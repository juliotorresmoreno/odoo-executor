package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juliotorresmoreno/odoo-executor/backup"
	"github.com/juliotorresmoreno/odoo-executor/config"
)

type Handler struct {
	backup *backup.OdooBackup
}

func ConfigureHandler() http.Handler {
	config := config.GetConfig()

	h := &Handler{
		backup: backup.NewOdooBackup(backup.OdooBackupConfig{
			OdooURL:        config.AdminURL,
			MasterPassword: config.AdminPassword,
			Namespace:      config.Namespace,
		}),
	}

	mux := mux.NewRouter()

	// Define your routes here
	mux.HandleFunc("/health", h.healthCheckHandler).Methods("GET")
	mux.HandleFunc("/backup/{id}", h.backupHandler).Methods("POST")
	mux.HandleFunc("/list", h.listHandler).Methods("GET")
	mux.HandleFunc("/download/{id}", h.downloadHandler).Methods("GET")

	return mux

}
