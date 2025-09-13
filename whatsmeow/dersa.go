package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/redis/go-redis/v9"
)

// parseInt safely converts string → int
func parseInt(s string) int {
	if num, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
		return num
	}
	return 0
}

// fetchRouteInfo scrapes one route and returns the report as a string
func fetchRouteInfo(routeID string) (string, error) {
	res, err := http.Get("https://semil.sp.gov.br/travessias/")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	fromLocation := doc.Find(fmt.Sprintf("#menu-travessia-a-%s", routeID)).Text()
	toLocation := doc.Find(fmt.Sprintf("#menu-travessia-b-%s", routeID)).Text()
	timeFrom := parseInt(doc.Find(fmt.Sprintf("#menu-travMinutosA-%s", routeID)).Text())
	timeTo := parseInt(doc.Find(fmt.Sprintf("#menu-travMinutosB-%s", routeID)).Text())
	vessels := parseInt(doc.Find(fmt.Sprintf("#menu-embarcacao-%s strong", routeID)).Text())
	conditions := doc.Find(fmt.Sprintf("#menu-tempoClima-%s", routeID)).AttrOr("title", "Unknown")
	// routeTitle := doc.Find(fmt.Sprintf("#menu-title-%s", routeID)).AttrOr("title", "Unknown")

	// Build a string instead of printing
	report := fmt.Sprintf(
		`=== EXTRAÍDO DE SEMIL-SP ===
Rota: SANTOS/GUARUJÁ
%s → %s: %d minutos
%s → %s: %d minutos
Nº de embarcações: %d
Clima/Tempo: %s
Atualizado: %s`,
		fromLocation, toLocation, timeFrom,
		toLocation, fromLocation, timeTo,
		vessels,
		conditions,
		time.Now().Format("15:04:05"),
	)

	return report, nil
}

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func getRouteInfoWithCache(routeID string, ctx context.Context) (string, error) {
	// Try cache first
	val, err := rdb.Get(ctx, routeID).Result()
	if err == nil {
		fmt.Println("✅ Got data from Redis cache")
		return val, nil
	}

	// Cache miss → scrape
	data, err := fetchRouteInfo(routeID)
	if err != nil {
		return "", err
	}

	// Store in Redis with 5-minute expiration
	err = rdb.Set(ctx, routeID, data, 5*time.Minute).Err()
	if err != nil {
		return "", err
	}

	fmt.Println("🆕 Scraped new data and cached in Redis")
	return data, nil
}

// func main() {
// 	routeID := "1951"
// 	report, err := fetchRouteInfo(routeID)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println(report)
// }
