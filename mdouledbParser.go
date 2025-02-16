package moduledbParser

import (
	"bytes"
	"github.com/mbock573/httpClientHelper"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"io"
	"log"
	"net/http"
	"strings"
)

type pflicht struct {
	name     string
	semester string
}

type wahl struct {
	name     string
	semester string
}

// Parses a single Module input is the specific URL for that course wich is
// retrievable with the courseParser package
// Url example https://moduldb.htwsaar.de/cgi-bin/moduldb-b?bkeys=avbb2a&lang=de
func Run(client *http.Client, courseUrl string) ([]pflicht, []wahl) {
	htmlResult, err := httpClientHelper.HttpGetRequest(client, courseUrl)
	if err != nil {
		log.Printf("moduledbParser: Error while http GET request: %v", err)
	}
	bodyBytes, err := io.ReadAll(htmlResult.Body)
	if err != nil {
		log.Printf("moduledbParser: Error while reading response body: %v", err)
	}
	pflicht, wahl := htmlParse(bodyBytes)
	return pflicht, wahl
}

// Parses the moduledb course page for a selected course
func htmlParse(htmlBody []byte) ([]pflicht, []wahl) {
	utf8Reader, err := charset.NewReader(bytes.NewReader(htmlBody), "")
	if err != nil {
		log.Printf("Fehler beim Konvertieren der Kodierung: %v", err)
		return nil, nil
	}

	// Parse das HTML
	doc, err := html.Parse(utf8Reader)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return nil, nil
	}

	var pflichtModule []pflicht
	var wahlModule []wahl
	var currentPflichth2 bool
	var currentWahlh2 bool
	var inTable bool

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		// Erkennung der Überschriften
		if n.Type == html.ElementNode && n.Data == "h2" {
			if getText(n) == "Pflichtmodule" {
				currentPflichth2 = true
				currentWahlh2 = false
				inTable = false
			} else if getText(n) == "Wahlmodule" {
				currentWahlh2 = true
				currentPflichth2 = false
				inTable = false
			}
		}

		// Tabellenverarbeitung
		if n.Type == html.ElementNode && n.Data == "table" {
			if currentPflichth2 || currentWahlh2 {
				inTable = true
			}
		}

		if inTable && n.Type == html.ElementNode && n.Data == "tr" {
			var modulName, semester string
			tdCount := 0

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.ElementNode && c.Data == "td" {
					tdCount++
					text := strings.TrimSpace(getText(c))

					switch tdCount {
					case 1: // Modulname
						modulName = text
					case 4: // Semester
						semester = text
					}
				}
			}

			if modulName != "" && semester != "" {
				if currentPflichth2 {
					pflichtModule = append(pflichtModule, pflicht{
						name:     modulName,
						semester: semester,
					})
				} else if currentWahlh2 {
					wahlModule = append(wahlModule, wahl{
						name:     modulName,
						semester: semester,
					})
				}
			}
		}

		// Rekursion durch Kinderknoten
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}

		// Reset nach Tabellenende
		if n.Type == html.ElementNode && n.Data == "table" {
			inTable = false
			currentPflichth2 = false
			currentWahlh2 = false
		}
	}

	traverse(doc)
	return pflichtModule, wahlModule
}

func getText(n *html.Node) string {
	var text strings.Builder
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(n)
	return strings.TrimSpace(text.String())
}
