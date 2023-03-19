package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", HomeHandler)
	mux.HandleFunc("/artist/", ArtistPage)
	mux.HandleFunc("/search/", Results)
	mux.HandleFunc("/filter/", Filter)
	mux.Handle("/ui/static/", http.StripPrefix("/ui/static/", http.FileServer(http.Dir("ui/static"))))

	// fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static")})
	// mux.Handle("/ui/static", http.NotFoundHandler())
	// mux.Handle("/ui/static/", http.StripPrefix("ui/static", fileServer))

	log.Println("Запуск веб-сервера на http://localhost:8070/ ")
	err := http.ListenAndServe(":8070", mux)
	if err != nil {
		log.Fatal(err)
	}
}

// type neuteredFileSystem struct {
// 	fs http.FileSystem
// }

// func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
// 	f, err := nfs.fs.Open(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s, err := f.Stat()
// 	if s.IsDir() {
// 		index := filepath.Join(path, "index.html")
// 		if _, err := nfs.fs.Open(index); err != nil {
// 			closeErr := f.Close()
// 			if closeErr != nil {
// 				return nil, closeErr
// 			}

// 			return nil, err
// 		}
// 	}

// 	return f, nil
// }
