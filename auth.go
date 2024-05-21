package main

import (
    "context"
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
    tokenString, err := token.SignedString(mySigningKey)

    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome! Please login or register.")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    username := "testuser"
    role := "organizer"

    token, err := GenerateJWT(username, role)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
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
            w.WriteHeader(http.StatusUnauthorized)
            fmt.Fprintf(w, "Invalid token")
            return
        }

        next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "props", claims)))
    })
}

func EventInfoHandler(w http.ResponseWriter, r *http.Request) {
    props := r.Context().Value("props").(*Claims)
    var info string

    switch props.Role {
    case "organizer":
        info = "Organizer-level event details"
    case "participant":
        info = "Participant-level event details"
    case "vendor":
        info = "Vendor-level event details"
    default:
        info = "Unknown role"
    }

    fmt.Fprintf(w, info)
}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/login", LoginHandler)
    r.Handle("/eventinfo", ValidateMiddleware(http.HandlerFunc(EventInfoHandler)))

    http.Handle("/", r)

    log.Fatal(http.ListenAndServe(":8080", nil))
}