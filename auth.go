package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

var mySigningKey []byte

func init() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
    mySigningKey = []byte(os.Getenv("SECRET_KEY"))
}

type Claims struct {
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.StandardClaims
}

type ErrorResponse struct {
    ErrorCode    int    `json:"error_code"`
    ErrorMessage string `json:"error_message"`
}

func GenerateJWT(username, role string) (string, error) {
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: username,
        Role:     role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(mySigningKey)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome! Please login or register.")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    // In a real application, you'd extract username and role from the request
    username := "testuser"
    role := "organizer"

    token, err := GenerateJWT(username, role)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Error generating token")
        return
    }

    fmt.Fprintf(w, "Token: %s", token)
}

func ValidateMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")

        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return mySigningKey, nil
        })

        if err != nil || !token.Valid {
            respondWithError(w, http.StatusUnauthorized, "Invalid token")
            return
        }

        next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "props", claims)))
    })
}

func EventInfoHandler(w http.ResponseWriter, r *http.Request) {
    claims := r.Context().Value("props").(*Claims)
    info := fmt.Sprintf("%s-level event details", claims.Role)

    fmt.Fprintf(w, info)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    response := ErrorResponse{ErrorCode: code, ErrorMessage: message}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        log.Printf("Error marshalling error response: %s", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(jsonResponse)
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/login", LoginHandler)
    r.Handle("/eventinfo", ValidateMiddleware(http.HandlerFunc(EventInfoHandler)))

    http.Handle("/", r)

    log.Fatal(http.ListenAndServe(":8080", nil))
}