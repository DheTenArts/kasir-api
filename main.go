package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Categories struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	KAT string `json:"kat"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var kategori = []Categories{
	{ID: 1, Nama: "Mie Goreng 100", KAT: "Makanan", Harga: 3000, Stok: 25},
	{ID: 2, Nama: "Teh Sari Murni", KAT: "Minuman", Harga: 2500, Stok: 30},
	{ID: 3, Nama: "Saos ABC", KAT: "Bumbu", Harga: 12000, Stok: 20},
	{ID: 4, Nama: "Teh Gelas", KAT: "Minuman", Harga: 1000, Stok: 10},
	{ID: 5, Nama: "Sarden DHE", KAT: "Makanan", Harga: 15000, Stok: 35},
}

// ! buat function untuk ambil kategori berdasarkan id
func getCategoriesByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	// URL: /api/category/123 -> ID = 123
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak angka", http.StatusBadRequest)
		return
	}

	// Cari kategori dengan ID tersebut
	for _, k := range kategori {
		if k.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(k)
			return
		}
	}

	// Kalau tidak found
	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

// ! buat fucntion untuk update isi kategori berdasarkan id
// ? localhost:8080/api/category/1
/*
	! PUT
	! di postman pilih body lalu raw terus bagian json isi dibawah ini
	! {
	!	"nama" : "Frozen Food",
	!	"harga": 25000,
	!	"stok": 20
	! }
	! maka kategori id 1 akan terganti/terupdate isinya
*/
func updateCategory(w http.ResponseWriter, r *http.Request) {
	// get id dari request
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")

	// ganti int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak angka!", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateCategory Categories
	err = json.NewDecoder(r.Body).Decode(&updateCategory)
	if err != nil {
		http.Error(w, "Request Gagal!", http.StatusBadRequest)
		return
	}

	// loop kategori, cari id, ganti sesuai data dari request
	for i := range kategori {
		if kategori[i].ID == id {
			updateCategory.ID = id
			kategori[i] = updateCategory

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateCategory)
			return
		}
	}

	http.Error(w, "Kategori belum ada", http.StatusNotFound)
}

// ! buat fucntion untuk delete kategori berdasarkan id
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	// get id
	idStr := strings.TrimPrefix(r.URL.Path, "/api/category/")

	// ganti id int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID tidak angka!", http.StatusBadRequest)
		return
	}

	// loop kategori cari ID, dapet index yang mau dihapus
	for i, p := range kategori {
		if p.ID == id {
			// bikin slice baru dengan data sebelum dan sesudah index
			kategori = append(kategori[:i], kategori[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "kategori belum ada", http.StatusNotFound)
}

func main() {

	// ? GET localhost:8080/api/category/{id}
	// ? PUT localhost:8080/api/category/{id}
	// ? DELETE localhost:8080/api/category/{id}
	http.HandleFunc("/api/category/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getCategoriesByID(w, r)
		} else if r.Method == "PUT" {
			updateCategory(w, r)
		} else if r.Method == "DELETE" {
			deleteCategory(w, r)
		}

	})

	// ! Method GET dan POST
	/*
		! POST
		! di postman pilih body lalu raw terus bagian json isi dibawah ini
		! {
		!	"nama" : "Frozen Food",
		!	"harga": 25000,
		!	"stok": 20
		! }
		! maka akan nampil kategori baru dengan id 6
	*/
	// ! localhost:8080/api/category
	http.HandleFunc("/api/category", func(w http.ResponseWriter, r *http.Request) {
		// ? Ambil kategori dari var kategori
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(kategori)
		} else if r.Method == "POST" {

			// ? buat variable baru
			var CategoriesBaru Categories
			err := json.NewDecoder(r.Body).Decode(&CategoriesBaru)
			if err != nil {
				http.Error(w, "Request gagal!", http.StatusBadRequest)
				return
			}

			// ? masukkin data ke dalam variable kategori
			CategoriesBaru.ID = len(kategori) + 1
			kategori = append(kategori, CategoriesBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(CategoriesBaru)
		}
	})

	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})
	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server")
	}
}
