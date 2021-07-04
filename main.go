package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jacobsa/go-serial/serial"
)

const (
	serialPort    = "/dev/tty.usbserial-220"
	listenAddr    = "10.69.69.250:8000"
	commandString = "swi0%s\n"
)

func main() {
	options := serial.OpenOptions{
		PortName:        serialPort,
		BaudRate:        19200,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}

	r := mux.NewRouter()
	r.HandleFunc("/kvm/{port}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		portNumber := vars["port"]

		if _, err := strconv.ParseInt(portNumber, 10, 8); err != nil {
			log.Fatalf("invalid port: %q", portNumber)
		}

		s, err := serial.Open(options)
		if err != nil {
			log.Println("No device connected at", serialPort)
			os.Exit(0)
		}
		defer s.Close()

		_, err = s.Write([]byte("Open\n"))
		if err != nil {
			log.Println("Unable to write to connected device")
			os.Exit(0)
		}

		_, err = s.Write([]byte(fmt.Sprintf(commandString, portNumber)))
		if err != nil {
			log.Println("Unable to write to connected device")
			os.Exit(0)
		}
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
