package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Structure pour un produit
type Product struct {
	Name        string
	Price       string
	Image       string
	ReducPrice  string
	ID          int
	Description string
	Stock       int
	Size        string
}

// Liste de produits
var products = []Product{
	{Name: "Palace Pull à Capuche Unisexe Chasseur", Price: "145€", Image: "/static/img/products/16A.webp", ID: 1, Description: "Un pull à capuche confortable pour un style décontracté.", Stock: 10, Size: "M"},
	{Name: "Sweat Nike Sportswear", Price: "89€", Image: "/static/img/products/18A.webp", ID: 2, Description: "Un sweat confortable pour un look décontracté.", Stock: 15, Size: "L"},
	{Name: "Veste Adidas Originals", Price: "120€", Image: "/static/img/products/19A.webp", ID: 3, Description: "Une veste élégante pour les journées fraîches.", Stock: 5, Size: "M"},
	{Name: "Veste The North Face", Price: "200€", Image: "/static/img/products/21A.webp", ID: 4, Description: "Une veste robuste pour les aventures en plein air.", Stock: 8, Size: "L"},
	{Name: "Parka Carhartt WIP", Price: "250€", ReducPrice: "150€", Image: "/static/img/products/22A.webp", ID: 5, Description: "Une parka chaude et stylée pour l'hiver.", Stock: 12, Size: "XL"},
	{Name: "Blouson Levi's Trucker ", Price: "110€", Image: "/static/img/products/33B.webp", ID: 6, Description: "Un blouson en denim intemporel.", Stock: 20, Size: "M"},
}

func main() {

	temp, errtemp := template.ParseGlob("./assets/temp/*.html")
	if errtemp != nil {
		fmt.Println(errtemp)
		os.Exit(1)
	}

	http.HandleFunc("/assets/temp/index", func(w http.ResponseWriter, r *http.Request) {

		temp.ExecuteTemplate(w, "index", products)
	})

	http.HandleFunc("/product/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		idStr := strings.TrimPrefix(path, "/product/")

		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Erreur : ID invalide"))
			return
		}

		var foundProduct *Product
		for i := range products {
			if products[i].ID == id {
				foundProduct = &products[i]
				break
			}
		}

		if foundProduct == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Erreur 404 : Produit introuvable"))
			return
		}

		temp.ExecuteTemplate(w, "product", foundProduct)
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		temp.ExecuteTemplate(w, "add", nil)
	})

	http.HandleFunc("/add/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseMultipartForm(10 << 20); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Erreur: impossible de lire le formulaire"))
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		price := strings.TrimSpace(r.FormValue("price"))
		reduc := strings.TrimSpace(r.FormValue("reduc"))
		description := strings.TrimSpace(r.FormValue("description"))
		stockStr := strings.TrimSpace(r.FormValue("stock"))
		size := strings.TrimSpace(r.FormValue("size"))

		if name == "" || price == "" || description == "" || stockStr == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Erreur: tous les champs obligatoires ne sont pas remplis"))
			return
		}

		stock, err := strconv.Atoi(stockStr)
		if err != nil || stock < 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Erreur: stock invalide"))
			return
		}

		newID := 1
		if len(products) > 0 {
			newID = products[len(products)-1].ID + 1
		}

		defaultImage := "/static/img/products/00958058-veste-lazare-gris-moyen_5_1.avif"

		imageURL := ""
		file, header, fileErr := r.FormFile("image")
		if fileErr == nil && header != nil && header.Filename != "" {
			ext := strings.ToLower(filepath.Ext(header.Filename))
			switch ext {
			case ".png", ".jpg", ".jpeg", ".webp", ".avif", ".gif":
				filename := fmt.Sprintf("upload_%d%s", time.Now().UnixNano(), ext)
				outPath := filepath.Join("assets", "img", "products", filename)

				if out, err := os.Create(outPath); err == nil {
					defer out.Close()
					defer file.Close()
					if _, errCopy := io.Copy(out, file); errCopy == nil {
						imageURL = "/static/img/products/" + filename
					}
				}
			}
		}
		if imageURL == "" {
			imageURL = defaultImage
		}

		newProduct := Product{
			Name:        name,
			Price:       price,
			ReducPrice:  reduc,
			Image:       imageURL,
			ID:          newID,
			Description: description,
			Stock:       stock,
			Size:        size,
		}

		products = append(products, newProduct)

		http.Redirect(w, r, "/product/"+strconv.Itoa(newID), http.StatusSeeOther)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() != "/" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404"))
			return
		}
		temp.ExecuteTemplate(w, "index", products)
	})

	chemin, _ := os.Getwd()
	fmt.Println(chemin)
	fileserver := http.FileServer(http.Dir(chemin + "/assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	http.ListenAndServe("localhost:8080", nil)

}
