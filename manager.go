package main

import (
	"errors"
	"log"
	"path/filepath"
	"sync"
)

const (
	UPLOAD_CONTENT_MESSAGE = "Upload content"
	UPLOAD_LOCAL_CONTENT   = "Upload local content to remote"

	SYNCHRONIZE_CONTENT_MESSAGE = "Synchronize content"
	DOWNLOAD_REMOTE_CONTENT     = "Download remote content."
	CONTENT_NEWER               = "The local content is newer than remote."
)

type Manager struct {
	Config Config
}

func New(config Config) Manager {
	return Manager{
		Config: config,
	}
}

func (m Manager) UploadContent() error {

	if !exists(m.Config.Remote) {
		return errors.New("Failed to synchronize content, remote directory does not exists.")
	}

	if !exists(m.Config.Local) {
		return errors.New("Failed to synchronize content, local directory does not exists.")
	}

	log.Println("Start upload content.")
	pushToast(UPLOAD_CONTENT_MESSAGE, UPLOAD_LOCAL_CONTENT)
	files := getFilesFromPattern(m.Config.Local, m.Config.Pattern)

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		from := filepath.Join(m.Config.Local, file)
		to := filepath.Join(m.Config.Remote, file)
		log.Println("--> Copy " + from + " to " + to)
		go func() {
			copyFile(from, to)
			wg.Done()
		}()
	}
	wg.Wait()

	writeHistory(m.Config.Local, files...)
	pushHistory(m.Config)

	return nil
}

func (m Manager) SynchonizeContent() error {
	log.Println("Check if synchronization is needed")
	if !exists(m.Config.Remote) {
		return errors.New("Failed to synchronize content, remote directory does not exists.")
	}

	if !exists(m.Config.Local) {
		return errors.New("Failed to synchronize content, local directory does not exists.")
	}

	remoteLastUpdate, _ := getRemoteLastUpdate(m.Config)
	localLastUpadate, _ := getLocalLastUpdate(m.Config)

	if remoteLastUpdate == 0 {
		log.Println("Nothing to synchonize.")
		return nil
	}

	if localLastUpadate == 0 {
		m.doSynchronize()
		return nil
	}

	if localLastUpadate < remoteLastUpdate {
		m.doSynchronize()
	} else {
		log.Println("Local content is older than remote, nothing to synchronize.")
		pushToast(SYNCHRONIZE_CONTENT_MESSAGE, CONTENT_NEWER)
	}

	return nil
}

func (m Manager) doSynchronize() {
	log.Println("Start download remote content.")
	pushToast(SYNCHRONIZE_CONTENT_MESSAGE, DOWNLOAD_REMOTE_CONTENT)
	files := getFilesFromPattern(m.Config.Remote, m.Config.Pattern)

	var wg sync.WaitGroup
	wg.Add(len(files))
	for _, file := range files {
		from := filepath.Join(m.Config.Remote, file)
		to := filepath.Join(m.Config.Local, file)
		log.Println("--> Copy " + from + " to " + to)
		go func() {
			copyFile(from, to)
			wg.Done()
		}()
	}
	wg.Wait()
	pullHistory(m.Config)
}
