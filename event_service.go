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
    http.HandleFunc("/event/", handleEventByID)

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
    w.Header().Set("Content-Type", "application/json")
    switch r.Method {
    case "POST":
        var event Event
        if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }

        mu.Lock()
        events[event.ID] = event
        mu.Unlock()

        w.WriteHeader(http.StatusCreated)
        if err := json.NewEncoder(w).Encode(event); err != nil {
            log.Printf("Error encoding event: %v\n", err)
            http.Error(w, "Failed to encode event", http.StatusInternalServerError)
        }
    case "GET":
        mu.RLock()
        allEvents := make([]Event, 0, len(events))
        for _, event := range events {
            allEvents = append(allEvents, event)
        }
        mu.RUnlock()

        if err := json.NewEncoder(w).Encode(allEvents); err != nil {
            log.Printf("Error encoding all events: %v\n", err)
            http.Error(w, "Failed to encode events list", http.StatusInternalServerError)
        }
    default:
        w.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func handleEventByID(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/event/"):]

    mu.RLock()
    event, ok := events[id]
    mu.RUnlock()

    w.Header().Set("Content-Type", "application/json")
    if !ok {
        http.NotFound(w, r)
        return
    }

    switch r.Method {
    case "GET":
        if err := json.NewEncoder(w).Encode(event); err != nil {
            log.Printf("Error encoding event: %v\n", err)
            http.Error(w, "Failed to encode the event", http.StatusInternalServerError)
        }
    case "PUT":
        var updatedEvent Event
        if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
            http.Error(w, "Invalid request body", http.StatusBadRequest)
            return
        }
        updatedEvent.ID = id

        mu.Lock()
        events[id] = updatedEvent
        mu.Unlock()

        if err := json.NewEncoder(w).Encode(updatedEvent); err != nil {
            log.Printf("Error encoding updated event: %v\n", err)
            http.Error(w, "Failed to encode updated event", http.StatusInternalServerError)
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