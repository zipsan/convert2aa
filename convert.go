package main

import (
	"fmt"
	"html/template"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"

	"github.com/nfnt/resize"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// get body content
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// load image data
		img, _, err := image.Decode(file)
		file.Close()
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		HTMLTags := "<div id=\"convert2aa\"><pre style=\"font: 10px/5px monospace; font-family: 'Courier New', 'Monospace';letter-spacing: -1px;\">"
		width, height := uint(float64(img.Bounds().Max.X)*0.1), uint(float64(img.Bounds().Max.Y)*0.1)
		resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
		bounds := resizedImg.Bounds()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := resizedImg.At(x, y).RGBA()
				colorTag := fmt.Sprintf(
					"#%s%s%s\n",
					fmt.Sprintf("%02x", uint8(r)),
					fmt.Sprintf("%02x", uint8(g)),
					fmt.Sprintf("%02x", uint8(b)),
				)
				spanTag := fmt.Sprintf(
					"<span style=\"color:%s\">■</span>", colorTag,
				)
				HTMLTags += spanTag
			}
			HTMLTags += "<br />"
		}
		HTMLTags += "</pre></div>"

		type Response struct {
			HTMLTag string `json:"html"`
		}
		type HTMLResponse struct {
			HTMLTag template.HTML `json:"html"`
		}
		response := HTMLResponse{HTMLTag: template.HTML(HTMLTags)}
		//json.NewEncoder(w).Encode(response)

		tmpl, err := template.ParseFiles("templates/result.html")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, response)
	} else {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}

func main() {
	fmt.Println("starting convert2aa")
	http.HandleFunc("/", handler)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
