package main

import (
	"fmt"
	"html/template"
	"image"
	"image/color"
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

		// image width and height resized the image
		width, height := uint(float64(img.Bounds().Max.X)*0.2), uint(float64(img.Bounds().Max.Y)*0.2)

		// resized image object
		resizedImg := resize.Resize(width, height, img, resize.Lanczos3)
		bounds := resizedImg.Bounds()

		// HTMLTag String definitions
		HeadString := "<div id=\"convert2aa\"><pre style=\"font: 10px/5px monospace; font-family: 'Courier New', 'Monospace';letter-spacing: -1px;\">"
		ColorString := "#%s%s%s"
		TagString := "<span style=\"color:%s\">#</span>"
		FootString := "</pre></div>"
		BrString := "<br />"

		byteHeadString := []byte(HeadString)
		byteFootString := []byte(FootString)
		byteBrString := []byte(BrString)

		// generate HTMLTag bytesarray object
		// faster than joining strings
		maxsize :=
			len(HeadString) + // header
				bounds.Min.X*bounds.Min.Y*(len(ColorString)+len(TagString)) + //tag + colortag
				bounds.Min.Y*len(BrString) + // br
				len(FootString) // footer
		HTMLTags := make([]byte, 0, maxsize)
		HTMLTags = append(HTMLTags, byteHeadString...)

		// convert image to AA
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := color.RGBAModel.Convert(resizedImg.At(x, y)).RGBA()
				byteSpanString := []byte(
					fmt.Sprintf(
						TagString,
						fmt.Sprintf(
							ColorString,
							fmt.Sprintf("%02x", uint8(r)),
							fmt.Sprintf("%02x", uint8(g)),
							fmt.Sprintf("%02x", uint8(b)),
						),
					),
				)
				HTMLTags = append(HTMLTags, byteSpanString...)
			}
			HTMLTags = append(HTMLTags, byteBrString...)
		}
		HTMLTags = append(HTMLTags, byteFootString...)

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
