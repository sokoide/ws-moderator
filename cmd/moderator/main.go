package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/minio/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/sokoide/ws-ai/pkg/claude"
	"github.com/sokoide/ws-ai/pkg/dalle"
	"github.com/sokoide/ws-ai/pkg/pdf"
)

// types
type options struct {
	port            int
	logLevel        string
	level           log.Level
	imageBasePath   string
	imageHostIPAddr string
}
type ClaudeComm struct {
	user string
	cin  chan claude.Request
	cout chan claude.Response
}

// globals
var o options = options{
	port:     80,
	logLevel: "INFO",
}

var db Database = Database{}

var claudes sync.Map

const staticImageDir = "images"

// functions
func parseArgs() {
	flag.IntVar(&o.port, "port", o.port, "Port to listen on")
	flag.StringVar(&o.logLevel, "logLevel", o.logLevel, "Log level")
	flag.StringVar(&o.imageBasePath, "imageBasePath", "./images", "Path to image base directory")
	flag.StringVar(&o.imageHostIPAddr, "imageHostIPAddr", "127.0.0.1", "Your IP address for React to download images from")
	flag.Parse()

	level, err := log.ParseLevel(o.logLevel)
	if err == nil {
		o.level = level
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			// Respond to preflight requests
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// /go/messages
func getMessageHandler(w http.ResponseWriter, r *http.Request) {
	msgid := r.URL.Query().Get("msgid")
	req := loadRequestsForMsgID(msgid)

	if req == nil {
		log.Errorf("Failed to get document for msgid=%s", msgid)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintln(w, req.Message.Data)
}

// /go/pdf
func pdfHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	user := r.URL.Query().Get("user")
	// email := r.URL.Query().Get("email")
	textID := r.URL.Query().Get("msgid_txt")
	imageID := r.URL.Query().Get("msgid_url")

	req := loadRequestsForMsgID(textID)
	if req == nil {
		log.Errorf("Failed to get document for textID=%s", textID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	longText := req.Message.Data

	req = loadRequestsForMsgID(imageID)
	if req == nil {
		log.Errorf("Failed to get document for imageID=%s", imageID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	comps := strings.Split(req.Message.Data, "/")
	if len(comps) < 2 {
		log.Errorf("Wrong image URL %s for imageID=%s", req.Message.Data, imageID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	emailDir := comps[len(comps)-2]
	imageFile := comps[len(comps)-1]
	staticImageDirFull, _ := filepath.Abs(staticImageDir)
	imagePath := path.Join(staticImageDirFull, emailDir, imageFile)

	log.Debugf("using image %s", imagePath)
	fileName := emailDir + ".pdf"
	outputPath := path.Join(staticImageDirFull, emailDir, fileName)
	log.Debugf("outputPath %s", outputPath)

	err := pdf.GeneratePdf(title, "Family Day 2024", imagePath, user, longText, outputPath)
	if err != nil {
		log.Errorf("Failed to generate a PDF. err: %v", err)
		http.Error(w, "Failed o generate a PDF", http.StatusInternalServerError)
		return
	}

	// Get file info to determine the size
	fileInfo, err := os.Stat(outputPath)
	if err != nil {
		log.Errorf("Failed to stat the generated PDF file %s. err: %v", outputPath, err)
		http.Error(w, "Failed to stat the generated PDF file", http.StatusInternalServerError)
		return
	}

	// Open the file
	file, err := os.Open(outputPath)
	if err != nil {
		log.Errorf("Generated PDF file %s not found. err: %v", outputPath, err)
		http.Error(w, "Generated PDF file not found", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// Copy the file content to the response
	_, err = file.Seek(0, 0) // Reset file pointer to the start
	if err != nil {
		log.Errorf("Failed to seek file. err: %v", outputPath, err)
		http.Error(w, "Failed to seek file", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, file) // Copy file content to response
	if err != nil {
		log.Errorf("Failed to serve file. err: %v", outputPath, err)
		http.Error(w, "Failed to serve file", http.StatusInternalServerError)
	}
}

// /go/completion
func completionHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	user := r.URL.Query().Get("user")
	employee := r.URL.Query().Get("employee")
	email := r.URL.Query().Get("email")
	textID := r.URL.Query().Get("msgid_txt")
	imageID := r.URL.Query().Get("msgid_url")

	err := storeCompletion(title, user, employee, email, textID, imageID)
	if err != nil {
		log.Errorf("Failed to store completion for %s", email)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func startModerator() {
	log.Info("Starting WebSocket Moderator Server...")

	var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections by returning true
			return true
			// Alternatively,
			// return r.Header.Get("Origin") == "http://allowed-origin.com"
		}}

	// "/" for react
	staticDir := filepath.Join("gui", "out")
	staticDirFull, _ := filepath.Abs(staticDir)

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", http.StripPrefix("/", fs))
	log.Infof("http://%s/ is open for browsers, serving from %s", o.imageHostIPAddr, staticDirFull)

	// "/images" for images
	staticImageDirFull, _ := filepath.Abs(staticImageDir)

	fsImage := http.FileServer(http.Dir(staticImageDir))

	// http.Handle("/images/", http.StripPrefix("/images/", fsImage))
	imageHandler := http.StripPrefix("/images/", fsImage)
	http.Handle("/images/", withCORS(imageHandler))

	log.Infof("http://%s/images is open for images, serving from %s", o.imageHostIPAddr, staticImageDirFull)

	// get message
	http.Handle("/go/message", withCORS(http.HandlerFunc(getMessageHandler)))

	// pdf generation
	http.Handle("/go/pdf", withCORS(http.HandlerFunc(pdfHandler)))

	// completion
	http.Handle("/go/completion", withCORS(http.HandlerFunc(completionHandler)))

	// moderator websocket
	http.HandleFunc("/go/moderator", func(writer http.ResponseWriter, request *http.Request) {
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Info(err)
			return
		}
		defer ws.Close()

		moderatorID := uuid.NewString()
		msg := makeModRequestJsonBytes("", "bot", "", "txt", fmt.Sprintf("moderatorID: %s", moderatorID), true, false)
		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Error(err)
			return
		}

		monitor := &NewRequestMonitor{
			ID:   moderatorID,
			Conn: ws,
		}
		db.register(monitor)
		defer db.unregister(monitor)

		log.Infof("[%s] New Connection Established", moderatorID)

		// ping
		go pingModerator(ws, moderatorID)

		// receive loop
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] ReadMessage failed, %v", moderatorID, err)
				break
			}
			log.Infof("[%s] Received Message: %s", moderatorID, string(message))

			// convert message into ModRequest
			var req ModRequest
			err = json.Unmarshal(message, &req)
			if err != nil {
				log.Errorf("[%s] failed to parse %s", moderatorID, string(message))
				continue
			}
			if req.Message.Kind == "system" {
				switch req.Message.Data {
				case "":
					// return all non-moderated documents
					reqs := loadRequests(false)
					for _, req := range reqs {
						msg := makeModRequestJsonBytes(req.ID, req.ClientID, req.UserEmail, req.Message.Kind, req.Message.Data, req.Approved, req.Moderated)
						log.Debugf("sending %+v", req)
						err = ws.WriteMessage(websocket.TextMessage, msg)
						if err != nil {
							// if disconnected, it comes here
							log.Warnf("[%s] WriteMessage failed, %v", moderatorID, err)
							break
						}
					}
					continue
				case "approve":
					updatedRequest := updateRequest(req.ID, true, true)
					if updatedRequest == nil {
						log.Errorf("failed to approve %v", req.ID)
						continue
					}
					log.Infof("approved %v", req.ID)
					continue
				case "deny":
					updatedRequest := updateRequest(req.ID, false, true)
					if updatedRequest == nil {
						log.Errorf("failed to deny %v", req.ID)
						continue
					}
					log.Infof("denied %v", req.ID)

					// record a message for client
					storeRequest("bot", updatedRequest.UserEmail,
						fmt.Sprintf("Moderator denied the response from AI: %s", updatedRequest.ID),
						"txt", true, true)
					continue
				default:
					log.Warnf("[%s] unknown Message.Data %s. Continuing...", moderatorID, req.Message.Data)
					continue
				}
			}
		}

	})

	// chat (user) websocket
	http.HandleFunc("/go/chat", func(writer http.ResponseWriter, request *http.Request) {
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Info(err)
			return
		}
		defer ws.Close()

		clientID := uuid.NewString()
		msg := makeModRequestJsonBytes("", "bot", "", "txt", fmt.Sprintf("clientID: %s", clientID), true, false)
		err = ws.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Error(err)
			return
		}

		monitor := &DatabaseMonitor{
			ID:   clientID,
			Conn: ws,
		}
		db.register(monitor)
		defer db.unregister(monitor)

		log.Infof("[%s] New Connection Established", clientID)

		// receive loop
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] ReadMessage failed, %v", clientID, err)
				break
			}
			log.Infof("[%s] Received Message: %s", clientID, string(message))

			// convert message into ModRequest
			var req ModRequest
			err = json.Unmarshal(message, &req)
			if err != nil {
				log.Errorf("[%s] failed to parse %s", clientID, string(message))
				continue
			}
			if req.Message.Kind == "system" && req.Message.Data == "" {
				// return all approved documents for the user
				reqs := loadRequestsForUserEmail(req.UserEmail)
				for _, req := range reqs {
					msg := makeModRequestJsonBytes(req.ID, req.ClientID, req.UserEmail, req.Message.Kind, req.Message.Data, req.Approved, req.Moderated)
					log.Debugf("sending %+v", req)
					err = ws.WriteMessage(websocket.TextMessage, msg)
					if err != nil {
						// if disconnected, it comes here
						log.Warnf("[%s] WriteMessage failed, %v", clientID, err)
						break
					}
				}
				continue
			}

			// save it
			savedRequestID := storeRequest(clientID, req.UserEmail, req.Message.Data, req.Message.Kind, true, true)

			// send a message to client.
			msg := makeModRequestJsonBytes("", "bot", "", "txt",
				fmt.Sprintf("Moderating & generating %s", savedRequestID),
				true, false)
			err = ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				// if disconnected, it comes here
				log.Warnf("[%s] WriteMessage failed, %v", clientID, err)
				break
			}

			// call AI
			if strings.HasPrefix(req.Message.Data, "/imagine") {
				// image generation
				prompt := req.Message.Data[len("/imagine"):]
				prompt = strings.TrimSpace(prompt)

				dalleReq := dalle.DalleRequest{
					Model:  "dall-e-3",
					Prompt: prompt,
					Size:   "1024x1024",
					N:      1,
				}
				imageURL, err := dalle.GenerateImage(dalleReq)
				if err == nil {
					log.Infof("[%s] Image generated at %s", clientID, imageURL)
					log.Debugf("[%s] downloading %s", clientID, imageURL)
					filename, err := downloadFile(imageURL, o.imageBasePath, req.UserEmail)
					if err == nil {
						url := "http://" + path.Join(o.imageHostIPAddr, o.imageBasePath, convertEmailToPath(req.UserEmail), filename)
						log.Infof("[%s] %s downloaded", clientID, url)
						storeRequest("bot", req.UserEmail, url, "url", false, false)
					} else {
						log.Errorf("[%s] Failed to download %s, err: %v", clientID, imageURL, err)
						storeRequest("bot", req.UserEmail, "Failed to download the generated image. Please retry", "txt", true, true)
					}
				} else {
					log.Errorf("[%s] Failed to generate an image: %v", clientID, err)
					storeRequest("bot", req.UserEmail, fmt.Sprintf("Failed to get an answer from AI. %v", err), "txt", true, true)
				}
			} else {
				// Text generation
				var c *ClaudeComm

				if value, ok := claudes.Load(req.UserEmail); ok {
					c = value.(*ClaudeComm)
				} else {
					log.Infof("%s not available, making...", req.UserEmail)
					c = &ClaudeComm{
						user: req.UserEmail,
						cin:  make(chan claude.Request, 1),
						cout: make(chan claude.Response, 1),
					}
					claudes.Store(req.UserEmail, c)
					go claude.StartConversation(c.user, c.cin, c.cout)
					defer func(c *ClaudeComm) {
						// remove c from claudes map
						claudes.Delete(c.user)
						close(c.cin)
					}(c)
				}
				log.Debugf("[%s] sending to c.cin", clientID)
				c.cin <- claude.Request{Prompt: req.Message.Data}
				res, ok := <-c.cout
				if ok && res.Succeeded {
					storeRequest("bot", req.UserEmail,
						res.Text,
						"txt", false, false)
				} else {
					log.Errorf("[%s] failed to receive from AI %v", clientID, res)
					storeRequest("bot", req.UserEmail,
						fmt.Sprintf("Failed to get an answer from AI. %v.", res.Text),
						"txt", true, true)
				}
			}
		}
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", o.port), nil))
}

// main
func main() {
	parseArgs()
	log.SetLevel(o.level)

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	startModerator()
}
