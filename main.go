package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/juliotorresmoreno/odoo-executor/backup"
)

func main() {
	backup := flag.String("backup", "", "Ejecutar backup")
	list := flag.Bool("list", false, "Listar backups")
	download := flag.String("download", "", "Descargar backup especificado")

	flag.Parse()

	if *backup != "" {
		runBackup(*backup)
	} else if *list {
		runList()
	} else if *download != "" {
		runDownload(*download)
	} else {
		log.Println("No se indicó ningún parámetro válido.")
	}
}

func runBackup(id string) {
	if id == "" {
		json.NewEncoder(os.Stdout).Encode([]interface{}{})
		return
	}

	bck := backup.NewOdooBackup(backup.OdooBackupConfig{
		OdooURL:        os.Getenv("ADMIN_URL"),
		MasterPassword: os.Getenv("ADMIN_PASSWORD"),
		Namespace:      os.Getenv("NAMESPACE"),
	})

	backupFile, err := bck.OdooDatabase(id)
	if err != nil {
		json.NewEncoder(os.Stdout).Encode([]interface{}{})
		return
	}

	log.Println(backupFile)
}

func runList() {
	bck := backup.NewOdooBackup(backup.OdooBackupConfig{
		OdooURL:        os.Getenv("ADMIN_URL"),
		MasterPassword: os.Getenv("ADMIN_PASSWORD"),
		Namespace:      os.Getenv("NAMESPACE"),
	})

	backups, err := bck.ListBackups()
	if err != nil {
		json.NewEncoder(os.Stdout).Encode([]interface{}{})
		return
	}

	json.NewEncoder(os.Stdout).Encode(backups)
}

func runDownload(id string) {
	if id == "" {
		return
	}

	bck := backup.NewOdooBackup(backup.OdooBackupConfig{
		OdooURL:        os.Getenv("ADMIN_URL"),
		MasterPassword: os.Getenv("ADMIN_PASSWORD"),
		Namespace:      os.Getenv("NAMESPACE"),
	})

	data, err := bck.DownloadBackup(id)
	if err != nil {
		return
	}

	os.Stdout.Write(data)
}
