package backup

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type OdooBackupConfig struct {
	OdooURL        string
	MasterPassword string
	Namespace      string
}

type OdooBackup struct {
	OdooURL        string
	MasterPassword string
	Namespace      string
	BackupFormat   string
	OutputDir      string
}

func NewOdooBackup(config OdooBackupConfig) *OdooBackup {
	if config.OdooURL == "" || config.MasterPassword == "" {
		log.Fatal("OdooURL and MasterPassword must be set")
	}

	if config.Namespace == "" {
		config.Namespace = "default"
	}

	return &OdooBackup{
		OdooURL:        config.OdooURL,
		MasterPassword: config.MasterPassword,
		Namespace:      config.Namespace,
		BackupFormat:   "zip",
		OutputDir:      "/data",
	}
}

func (o *OdooBackup) OdooDatabase(dbName string) (string, error) {
	if o.OdooURL == "" || o.MasterPassword == "" || dbName == "" {
		return "", fmt.Errorf("OdooURL, MasterPassword y dbName no pueden estar vacíos")
	}

	backupEndpoint := fmt.Sprintf("%s/web/database/backup", o.OdooURL)

	form := fmt.Sprintf("master_pwd=%s&name=%s&backup_format=%s", o.MasterPassword, dbName, o.BackupFormat)
	log.Printf("Iniciando backup para la DB '%s' en %s", dbName, o.OdooURL)

	req, err := http.NewRequest("POST", backupEndpoint, bytes.NewBufferString(form))
	if err != nil {
		return "", fmt.Errorf("error al crear la petición HTTP: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/octet-stream")

	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("no se pudo hacer el backup de la base de datos '%s': %w", dbName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("error en la respuesta del servidor: %s - %s", resp.Status, string(bodyBytes))
	}

	backupFileName := fmt.Sprintf("%s_%s.%s", dbName, time.Now().Format("20060102_150405"), o.BackupFormat)
	fullPath := filepath.Join(o.OutputDir, backupFileName)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("error al crear directorio '%s': %w", filepath.Dir(fullPath), err)
	}

	outFile, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("error al crear archivo '%s': %w", fullPath, err)
	}
	defer outFile.Close()

	n, err := io.Copy(outFile, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error al guardar backup: %w", err)
	}

	log.Printf("Backup exitoso: '%s' (%d bytes)", fullPath, n)
	return fullPath, nil
}

type Backup struct {
	Name      string `json:"name"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"createdAt"`
}

func (o *OdooBackup) ListBackups() ([]Backup, error) {
	var backups = make([]Backup, 0)
	fullPath := o.OutputDir
	files, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error al leer el directorio de backups: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("error al obtener información del archivo '%s': %w", file.Name(), err)
		}
		backups = append(backups, Backup{
			Name:      file.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime().Format(time.RFC3339),
		})
	}

	return backups, nil
}

func (o *OdooBackup) DownloadBackup(fileName string) (io.ReadCloser, error) {
	fullPath := filepath.Join(o.OutputDir, fileName)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error al obtener información del archivo '%s': %w", fileName, err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("el archivo '%s' es un directorio, no se puede descargar", fileName)
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo '%s': %w", fileName, err)
	}

	return file, nil
}
