package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Content struct {
	Text       string  `json:"text"`
	FontSize   float64 `json:"font_size"`
	Rect       [4]int  `json:"rect"`
	LineHeight float64 `json:"line_height"`
	Color      string  `json:"color"`
}

type GenerateRequest struct {
	ImageURL string    `json:"image_url"`
	Contents []Content `json:"contents"`
}

func DrawTextWithLineHeight(img draw.Image, rect image.Rectangle, text string, fontFace font.Face, lineHeight float64, color color.Color, maxWidth int) {
	textWidth := rect.Dx()
	textHeight := rect.Dy()
	lines := strings.Split(WrapText(text, fontFace, maxWidth), "\n")
	var totalHeight int
	lineHeights := []int{}

	for _, line := range lines {
		bounds, _ := font.BoundString(fontFace, line)
		lineHeightPx := (bounds.Max.Y - bounds.Min.Y).Ceil()
		lineHeights = append(lineHeights, lineHeightPx)
		totalHeight += int(float64(lineHeightPx) * lineHeight)
	}

	currentY := rect.Min.Y + (textHeight-totalHeight)/2

	for i, line := range lines {
		bounds, _ := font.BoundString(fontFace, line)
		lineWidth := (bounds.Max.X - bounds.Min.X).Ceil()
		lineX := rect.Min.X + (textWidth-lineWidth)/2
		dot := fixed.P(lineX, currentY+lineHeights[i])
		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color),
			Face: fontFace,
			Dot:  dot,
		}
		drawer.DrawString(line)
		currentY += int(float64(lineHeights[i]) * lineHeight)
	}
}

func hexToColor(hex string) color.Color {
	if strings.HasPrefix(hex, "#") {
		hex = hex[1:]
	}

	var r, g, b uint8
	switch len(hex) {
	case 3:
		fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r = r * 17
		g = g * 17
		b = b * 17
	case 6:
		fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	default:
		return color.White
	}
	return color.RGBA{R: r, G: g, B: b, A: 255}
}

func WrapText(text string, fontFace font.Face, maxWidth int) string {
	if len(text) == 0 {
		return ""
	}

	var wrappedText string
	words := strings.Fields(text)
	currentLine := ""

	for _, word := range words {
		if currentLine == "" {
			currentLine = word
		} else {
			testLine := currentLine + " " + word
			bounds, _ := font.BoundString(fontFace, testLine)
			if (bounds.Max.X - bounds.Min.X).Ceil() > maxWidth {
				wrappedText += currentLine + "\n"
				currentLine = word
			} else {
				currentLine = testLine
			}
		}
	}
	wrappedText += currentLine
	return wrappedText
}

func GenerateImage(req GenerateRequest) ([]byte, error) {
	resp, err := http.Get(req.ImageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch image: %v", err)
	}
	defer resp.Body.Close()

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		return nil, fmt.Errorf("URL does not point to a valid image")
	}

	img, err := imaging.Decode(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	rgba := imaging.Clone(img)

	fontBytes, err := os.ReadFile("./fonts/Quicksand-Regular.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to load font: %v", err)
	}
	fontParsed, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %v", err)
	}

	for _, content := range req.Contents {
		fontFace := truetype.NewFace(fontParsed, &truetype.Options{Size: content.FontSize})
		rect := image.Rect(content.Rect[0], content.Rect[1], content.Rect[2], content.Rect[3])
		textColor := hexToColor(content.Color)
		DrawTextWithLineHeight(rgba, rect, content.Text, fontFace, content.LineHeight, textColor, rect.Dx())
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, rgba)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	return buf.Bytes(), nil
}

func handleGenerateImage(w http.ResponseWriter, r *http.Request) {
	var req GenerateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON body: %v", err), http.StatusBadRequest)
		return
	}

	if req.ImageURL == "" || !strings.HasPrefix(req.ImageURL, "http") {
		http.Error(w, "Invalid or missing image_url", http.StatusBadRequest)
		return
	}

	for _, content := range req.Contents {
		if content.FontSize <= 0 || content.LineHeight <= 0 {
			http.Error(w, "Invalid font size or line height", http.StatusBadRequest)
			return
		}
	}

	imageBytes, err := GenerateImage(req)
	if err != nil {
		log.Printf("Error generating image: %v", err)
		http.Error(w, "Error generating image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "inline; filename=\"result.png\"")
	w.Write(imageBytes)
}

func recoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next(w, r)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message": "Image generated successfully!"}`)
	})

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
