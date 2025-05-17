package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/avinash/weatherportfolio/api"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
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
	tracer := otel.Tracer("weather-handler")
	ctx, span := tracer.Start(r.Context(), "indexHandler")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	ip := getClientIP(r)
	if ip == "127.0.0.1" || ip == "::1" {
		ip = os.Getenv("LOCAL_IP") // fallback to West Bengal IP
	}

	api.Logger.Info("User IP", zap.String("ip", ip))

	locationData, err := ipAPI.GetLocation(ctx, ip)

	if err != nil {
		api.Logger.Error("Failed to get location", zap.Error(err))
		locationData = &api.LocationData{Country: "Unknown", RegionName: "Unknown", Timezone: "Unknown"}
	}

	weatherData, err := api.GetWeather(ctx,
		fmt.Sprintf("%f", locationData.Lat),
		fmt.Sprintf("%f", locationData.Lon),
	)
	if err != nil {
		api.Logger.Error("☁️ Failed to get weather:", zap.Error(err))
	} else {
		api.Logger.Info("🌤️ Weather info",
			zap.String("description", weatherData.Description),
			zap.Int("code", weatherData.CurrentWeather.Weathercode),
		)
	}

	data := struct {
		Location *api.LocationData
		Weather  *api.WeatherResponse
	}{
		Location: locationData,
		Weather:  weatherData,
	}

	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		api.HandleError(w, "Template parsing error", err, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		api.HandleError(w, "Template executing error", err, http.StatusInternalServerError)
		return
	}
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	if err := godotenv.Load(); err != nil {
		api.Logger.Info("⚠️ No .env file found, using system environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx := context.Background()
	tp, err := api.InitTracer(ctx, "weather-portfolio")
	if err != nil {
		log.Fatal("❌ Failed to initialize tracer:", err)
	}
	defer func() {
		_ = tp.Shutdown(ctx)
	}()
	api.InitLogger()
	http.HandleFunc("/", indexHandler)

	log.Printf("🚀 Server running at: http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
