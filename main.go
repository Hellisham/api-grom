package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Product struct {
	ID          uint
	Name        string
	Description string
	Price       uint
	Tax         uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}

var db *gorm.DB

func connDB() {
	var err error
	conn := "host=localhost user=admin password=admin dbname=postgres port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&Product{})
	fmt.Println("Connected to database")
}

func getALLproductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var products []Product
	var productsResponse []ProductResponse
	if err := db.Select("Name", "Description", "Price").Find(&products).Error; err != nil {
		http.Error(w, "Error retrieving books", http.StatusInternalServerError)
		return
	}
	for _, product := range products {
		productsResponse = append(productsResponse, ProductResponse{
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})

	}
	json.NewEncoder(w).Encode(productsResponse)
}

func getProductById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Book ID", http.StatusBadRequest)
		return
	}
	var product Product
	if err := db.First(&product, id).Error; err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}
	productsResponse := ProductResponse{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	json.NewEncoder(w).Encode(productsResponse)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
	}
	if err := db.Create(&product).Error; err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
	}
	productResponse := ProductResponse{
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}
	json.NewEncoder(w).Encode(productResponse)
}
func main() {
	connDB()
	router := mux.NewRouter()
	router.HandleFunc("/products", getALLproductsHandler).Methods("GET")
	router.HandleFunc("/product/{id}", getProductById).Methods("GET")
	router.HandleFunc("/newproduct", createProductHandler).Methods("POST")
	port := "8080"
	fmt.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
