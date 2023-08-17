package main

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

const (
	PREFIX_COMMAND = "./"
	PROCESS_NAME   = "FactoryGame.exe"
	INTERVAL       = 5 * time.Second
)

// Executes handler function periodicly
// The ticker is interrupt if handler return true 
func checkPeriodic(handler func() bool, interrupt chan struct{}) {
	log.Println("Watch process " + PROCESS_NAME)
	ticker := time.NewTicker(INTERVAL)
	go func() {
		for {
			select {
			case <-ticker.C:
				if handler() {
					ticker.Stop()
					close(interrupt)
				}
			}
		}
	}()
}

// Starts game process
func startGame() {
	gameCmd := exec.Command(PREFIX_COMMAND + PROCESS_NAME)
	_, err := gameCmd.Output()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	config, err := loadConfiguration()
	if err != nil {
		panic(err)
	}

	m := New(config)

	running, err := isProcessRunning(PROCESS_NAME)
	if err != nil {
		log.Fatal("Failed to check process state")
		panic(err)
	}

	if running {
		log.Println("Game already start, wait to synchronize saves.")
	} else {
		if err := m.SynchonizeContent(); err != nil {
			log.Println(err.Error())
		}
		log.Println("Start game")
		startGame()
	}

	interrupt := make(chan struct{})
	checkPeriodic(func() bool {
		isRunning, err := isProcessRunning(PROCESS_NAME)
		if err != nil {
			panic(err)
		}

		if !isRunning {
			log.Println("Process " + PROCESS_NAME + " is not running start synchronization")
			if err := m.UploadContent(); err != nil {
				log.Fatalln(err)
			}
			return true
		}
		return false
	}, interrupt)

	<-interrupt
	fmt.Println("Press Enter Key to terminate program")
	fmt.Scanln()
}
