package main

import (
	"context"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/avinash/weatherportfolio/api"
	"github.com/joho/godotenv"
)

var ipAPI api.IPAPI

func getClientIP(r *http.Request) string {

	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ip := getClientIP(r)
	if ip == "127.0.0.1" || ip == "::1" {
		ip = os.Getenv("LOCAL_IP") // fallback to West Bengal my ip
	}

	log.Println("User IP:", ip)

	locationData, err := ipAPI.GetLocation(ctx, ip)
	if err != nil {
		log.Println("Failed to get location: ", err)
		locationData = &api.LocationData{Country: "Unknown", RegionName: "Unknown", Timezone: "Unknown"}
	}
	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, "Template parsing error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, locationData)
	if err != nil {
		http.Error(w, "Template executing error", http.StatusInternalServerError)
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback
	}

	http.HandleFunc("/", indexHandler)

	log.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
