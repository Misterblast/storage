package database

import (
	"database/sql"

	log "github.com/ghulammuzz/misterblast-storage/utils"
)

// My SQL Schema

type File struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	OriginURL string `json:"origin_url"`
	GCSURL    string `json:"gcs_url"`
	Status    string `json:"status"`
}

func CreateOrUpdate(tx *sql.Tx, key, originURL, gcsURL, status string) error {

	query := `
        INSERT INTO files_url (key, origin_url, gcs_url, status) 
        VALUES ($1, $2, $3, $4) 
        ON CONFLICT (key) 
        DO UPDATE SET origin_url = EXCLUDED.origin_url, gcs_url = EXCLUDED.gcs_url, status = EXCLUDED.status;
    `
	_, err := tx.Exec(query, key, originURL, gcsURL, status)
	if err != nil {
		log.Error("[Database.CreateOrUpdate] Failed to create or update file: %s", err.Error())
		return err
	}
	return nil
}

func GetGCSURL(originURL string) (string, string, error) {
	var gcsURL string
	var fileStatus string

	query := "SELECT gcs_url, status FROM files_url WHERE origin_url = $1"
	err := DBInstance.QueryRow(query, originURL).Scan(&gcsURL, &fileStatus)
	if err != nil {
		return "", "", err
	}
	return gcsURL, fileStatus, nil

}

func Delete(tx *sql.Tx, originURL string) error {
	query := "DELETE FROM files_url WHERE origin_url = $1"
	_, err := tx.Exec(query, originURL)
	if err != nil {
		return err
	}
	return nil

}

func UpdateGCSURL(originURL, gcsURL string) error {
	query := "UPDATE files_url SET gcs_url = $1 WHERE origin_url = $2"
	_, err := DBInstance.Exec(query, gcsURL, originURL)
	if err != nil {
		return err
	}
	return nil
}

func CheckKey(key string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM files_url WHERE key = $1)"
	var exists bool
	err := DBInstance.QueryRow(query, key).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
