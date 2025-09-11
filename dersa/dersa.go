package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	res, err := http.Get("https://semil.sp.gov.br/travessias/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Extract data using specific IDs
	routeID := "1951"
	fromLocation := doc.Find(fmt.Sprintf("#menu-travessia-a-%s", routeID)).Text()
	toLocation := doc.Find(fmt.Sprintf("#menu-travessia-b-%s", routeID)).Text()
	timeFrom := parseInt(doc.Find(fmt.Sprintf("#menu-travMinutosA-%s", routeID)).Text())
	timeTo := parseInt(doc.Find(fmt.Sprintf("#menu-travMinutosB-%s", routeID)).Text())
	vessels := parseInt(doc.Find(fmt.Sprintf("#menu-embarcacao-%s strong", routeID)).Text())
	conditions := doc.Find(fmt.Sprintf("#menu-tempoClima-%s", routeID)).AttrOr("title", "Unknown")
	routeTitle := doc.Find(fmt.Sprintf("#menu-title-%s", routeID)).AttrOr("title", "Unknown")

	// Display results
	fmt.Printf("=== DIRECT ID EXTRACTION ===\n")
	fmt.Printf("Route: %s\n", routeTitle)
	fmt.Printf("%s → %s: %d minutes\n", fromLocation, toLocation, timeFrom)
	fmt.Printf("%s → %s: %d minutes\n", toLocation, fromLocation, timeTo)
	fmt.Printf("Total time: %d minutes\n", timeFrom+timeTo)
	fmt.Printf("Vessels: %d\n", vessels)
	fmt.Printf("Weather: %s\n", conditions)
	fmt.Printf("Updated: %s\n", time.Now().Format("15:04:05"))
}

func parseInt(s string) int {
	if num, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
		return num
	}
	return 0
}
