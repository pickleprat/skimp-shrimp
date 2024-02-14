package _middleware

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func MiddlewareChain(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request, middleware ...func(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request)) {
	for _, mw := range middleware {
		mw(ctx, w, r)
	}
}

func Log(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	startTime := ctx["StartTime"].(time.Time)
	elapsedTime := time.Since(startTime)
	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] [%s] [%s]\n", formattedTime, r.Method, r.URL.Path, elapsedTime)
}

func Init(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	ctx["StartTime"] = time.Now()
}

func MwRandomTruett(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	quotes := []string{
		"Quality is remembered long after the price is forgotten.",
		"People don't care how much you know until they know how much you care.",
		"Be better before you get bigger.",
		"The more you give, the more you get.",
		"If you want a great business, you have to have great operators.",
		"It's easier to succeed than to fail.",
		"We live in a changing world, but we need to be reminded that the important things have not changed.",
		"Food is essential to life; therefore, make it good.",
		"It is more rewarding to watch money change the world than watch it accumulate.",
		"Adversity is the diamond dust heaven polishes its jewels with.",
	}
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)
	randomIndex := randomGenerator.Intn(len(quotes))
	randomQuote := quotes[randomIndex]
	ctx["RandomTruett"] = randomQuote
}

func ParseForm(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
}

func ParseMultipartForm(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
}

func MwGenMultipartString(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	randomBytes := make([]byte, 5)
	randString := hex.EncodeToString(randomBytes)
	ctx["MultipartString"] = randString
}

func Auth(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("SessionToken")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if cookie.Value != os.Getenv("ADMIN_SESSION_TOKEN") {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func MwIncludePNG(ctx map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
}