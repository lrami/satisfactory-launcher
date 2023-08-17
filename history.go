package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

const HISTORY = "history.json"

type History struct {
	LastUpdate int64
	Files      []string
}

func loadHistory(path string) (History, error) {
	historyPath := filepath.Join(path, HISTORY)
	content, err := ioutil.ReadFile(historyPath)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
		return History{}, err
	}

	var payload History
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
		return History{}, err
	}
	return payload, nil
}

func loadLocalHistory(config Config) (History, error) {
	historyPath := filepath.Join(config.Local, HISTORY)
	if !exists(historyPath) {
		return History{}, errors.New("Local history does not exists.")
	}
	return loadHistory(config.Local)
}

func loadRemoteHistory(config Config) (History, error) {
	historyPath := filepath.Join(config.Remote, HISTORY)
	if !exists(historyPath) {
		return History{}, errors.New("Remote history does not exists.")
	}
	return loadHistory(config.Remote)
}

func getLocalLastUpdate(config Config) (int64, error) {
	history, err := loadLocalHistory(config)
	if err != nil {
		return 0, errors.New("Unable to retrieve local history")
	}
	return history.LastUpdate, nil
}

func getRemoteLastUpdate(Config Config) (int64, error) {
	history, err := loadRemoteHistory(Config)
	if err != nil {
		return 0, errors.New("Unable to retrieve remote history")
	}
	return history.LastUpdate, nil
}

func pullHistory(config Config) {
	localHistoryPath := filepath.Join(config.Local, HISTORY)
	remoteHistoryPath := filepath.Join(config.Remote, HISTORY)
	if exists(localHistoryPath) {
		os.Remove(localHistoryPath)
	}
	log.Println("--> Pull remote history.json to local directory")
	copyFile(remoteHistoryPath, localHistoryPath)
}

func pushHistory(config Config) {
	localHistoryPath := filepath.Join(config.Local, HISTORY)
	remoteHistoryPath := filepath.Join(config.Remote, HISTORY)
	log.Println("--> Push local history.json to remote directory")
	copyFile(localHistoryPath, remoteHistoryPath)
}

func writeHistory(path string, files ...string) error {
	now := time.Now()
	timestamp := now.Unix()
	historyPath := filepath.Join(path, HISTORY)
	history := History{
		LastUpdate: timestamp,
		Files:      files,
	}

	raw, err := json.Marshal(history)
	if err != nil {
		return errors.New("Unable to encode history object.")
	}

	err = WriteFile(historyPath, string(raw))
	if err != nil {
		return errors.New("Unable to write history file.")
	}
	return nil
}
