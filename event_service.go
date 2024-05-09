package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "sync"
)

type Event struct {
    ID           string   `json:"id"`
    Name         string   `json:"name"`
    Participants []string `json:"participants"`
    Vendors      []string `json:"vendors"`
    Schedule     string   `json:"schedule"`
}

var (
    events = make(map[string]Event)
    mu     sync.RWMutex
)

func main() {
    http.HandleFunc("/events", handleEvents)
    http.HandleFunc("/event/", handleEvent)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Server starting on port %s\n", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func handleEvents(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        var event Event
        if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            log.Printf("Error decoding event: %v\n", err)
            return
        }

        mu.Lock()
        events[event.ID] = event
        mu.Unlock()

        w.WriteHeader(http.StatusCreated)
        if err := json.NewEncoder(w).Encode(event); err != nil {
            log.Printf("Error encoding event: %v\n", err)
        }
    case "GET":
        mu.RLock()
        var allEvents []Event
        for _, event := range events {
            allEvents = append(allEvents, event)
        }
        mu.RUnlock()

        if err := json.NewEncoder(w).Encode(allEvents); err != nil {
            log.Printf("Error encoding all events: %v\n", err)
        }
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/event/"):]

    mu.RLock()
    event, ok := events[id]
    mu.RUnlock()

    if !ok {
        http.NotFound(w, r)
        return
    }

    switch r.Method {
    case "GET":
        if err := json.NewEncoder(w).Encode(event); err != nil {
            log.Printf("Error encoding event: %v\n", err)
        }
    case "PUT":
        var updatedEvent Event
        if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            log.Printf("Error decoding updated event: %v\n", err)
            return
        }
        updatedEvent.ID = id

        mu.Lock()
        events[id] = updatedEvent
        mu.Unlock()

        if err := json.NewEncoder(w).Encode(updatedEvent); err != nil {
            log.Printf("Error encoding updated event: %v\n", err)
        }
    case "DELETE":
        mu.Lock()
        delete(events, id)
        mu.Unlock()

        w.WriteHeader(http.StatusNoContent)
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}