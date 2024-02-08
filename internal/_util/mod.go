package _util

import (
	"encoding/base64"
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
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

func ServeStaticFilesAndFavicon() {
	http.HandleFunc("/static/", ServeStaticFiles)
	http.HandleFunc("/favicon.ico", ServeFavicon)
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

func URLBuilder(baseURL string, queryParams ...string) string {
	if len(queryParams)%2 != 0 {
		return baseURL
	}
	urlParams := make(map[string]string)
	for i := 0; i < len(queryParams); i += 2 {
		key := queryParams[i]
		value := queryParams[i+1]
		urlParams[key] = value
	}
	url := baseURL
	if len(urlParams) > 0 {
		url += "?"
		for key, value := range urlParams {
			url += key + "=" + value + "&"
		}
		url = url[:len(url)-1]
	}
	return url
}