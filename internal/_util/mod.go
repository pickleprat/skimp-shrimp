package _util

import (
	"encoding/base64"
	"errors"
	"html/template"
	"math/rand"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func RenderTemplate(w http.ResponseWriter, tmplFile string, data map[string]interface{}) {
	tmpl, err := template.ParseFiles(tmplFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RenderErrorPage(w http.ResponseWriter, r *http.Request, err string) {
	http.Redirect(w, r, "/500?Err="+err, http.StatusSeeOther)
}

func ServeStaticFilesAndFavicon(mux *http.ServeMux) {
	mux.HandleFunc("GET /static/", ServeStaticFiles)
	mux.HandleFunc("GET /favicon.ico", ServeFavicon)
}

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	filePath := "favicon.ico"
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ConvertStringToUint(input string) (uint, error) {
	if input == "" {
		return 0, errors.New("input string is empty")
	}
	convertedUint, err := strconv.ParseUint(input, 10, 0)
	if err != nil {
		return 0, err
	}
	return uint(convertedUint), nil
}

func BytesToBase64String(input []byte) string {
	encodedString := base64.StdEncoding.EncodeToString(input)
	return encodedString
}

func URLBuilder(basePath string, params ...string) string {
    if len(params)%2 != 0 {
        panic("Invalid number of parameters. Must be even.")
    }

    if len(params) == 0 {
        return basePath
    }

    var sb strings.Builder
    sb.WriteString(basePath)
    sb.WriteString("?")

    for i := 0; i < len(params); i += 2 {
        if i > 0 {
            sb.WriteString("&")
        }
        sb.WriteString(params[i])
        sb.WriteString("=")
        sb.WriteString(params[i+1])
    }
    return sb.String()
}

func IsValidPhoneNumber(phone string) bool {
    pattern := `^\d{3}-\d{3}-\d{4}$`
    regex := regexp.MustCompile(pattern)
    return regex.MatchString(phone)
}

func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

func GenerateRandomToken(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    token := make([]byte, length)
    rand.Seed(time.Now().UnixNano())
    for i := range token {
        token[i] = charset[rand.Intn(len(charset))]
    }
    return string(token)
}