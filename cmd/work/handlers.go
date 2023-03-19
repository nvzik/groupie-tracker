package main

import (
	functions "groupie-tracker/pkg"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const (
	artist_groups string = "https://groupietrackers.herokuapp.com/api/artists"
	relation      string = "https://groupietrackers.herokuapp.com/api/relation"
)

type Error struct {
	CodeError        int
	ErrorDescription string
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	artists, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	Artists, err := functions.GetAllArtist(relation, artist_groups)
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	tml := artists.Execute(w, Artists)
	if tml != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func ArtistPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[:8] != "/artist/" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	id, err := strconv.Atoi(r.URL.Path[8:])
	if err != nil {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	if id > 52 || id <= 0 {
		Errors(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}
	artists, err := template.ParseFiles("./ui/html/artist-page.html")
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	Artist, err := functions.GetAllArtist(relation, artist_groups)
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	tml := artists.Execute(w, Artist[id-1])
	if tml != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
}

func Results(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search/" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	var res []functions.Artist
	query := r.FormValue("query")
	if !checktxt(query) && query != "" {
		res = []functions.Artist{}
	}
	Artists, err := functions.GetAllArtist(relation, artist_groups)
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	for _, v := range Artists {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(query)) && unique(res, v.Id) {
			res = append(res, v)
			continue
		}
		for _, members := range v.Members {
			if strings.Contains(strings.ToLower(members), strings.ToLower(query)) && unique(res, v.Id) {
				res = append(res, v)
				continue
			}
		}
		if strings.Contains(strconv.Itoa(v.CreationDate), query) && unique(res, v.Id) {
			res = append(res, v)
			continue
		}
		if strings.Contains(v.FirstAlbum, query) && unique(res, v.Id) {
			res = append(res, v)
			continue
		}
		for location := range v.DatesLocations {
			if strings.Contains(strings.ToLower(location), strings.ToLower(query)) && unique(res, v.Id) {
				res = append(res, v)
				continue
			}
		}
	}
	tmpl, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	result := tmpl.Execute(w, res)
	if result != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func Filter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		Errors(w, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	if r.URL.Path != "/filter/" {
		Errors(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		return
	}
	Artists, err := functions.GetAllArtist(relation, artist_groups)
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	CreationDateFrom, _ := strconv.Atoi(r.FormValue("CreationDateFrom"))
	CreationDateTo, _ := strconv.Atoi(r.FormValue("CreationDateTo"))
	FirstAlbumFrom, _ := strconv.Atoi(r.FormValue("FirstAlbumFrom"))
	FirstAlbumTo, _ := strconv.Atoi(r.FormValue("FirstAlbumTo"))
	Members := r.Form["members"]
	MembersInt := []int{}
	for _, c := range Members {
		ch, _ := strconv.Atoi(c)
		MembersInt = append(MembersInt, ch)
	}
	Locations := r.FormValue("searchLoc")
	var Res []functions.Artist

	for _, artist := range Artists {
		if Members == nil && Locations == "" {
			for i := CreationDateFrom; i <= CreationDateTo; i++ {
				for j := FirstAlbumFrom; j <= FirstAlbumTo; j++ {
					FirstAlbumInt, _ := strconv.Atoi(artist.FirstAlbum[6:])
					if i == artist.CreationDate && j == FirstAlbumInt {
						Res = append(Res, artist)
					}
				}
			}
		} else if Members != nil && Locations == "" {
			for i := CreationDateFrom; i <= CreationDateTo; i++ {
				for j := FirstAlbumFrom; j <= FirstAlbumTo; j++ {
					for _, member := range MembersInt {
						FirstAlbumInt, _ := strconv.Atoi(artist.FirstAlbum[6:])
						if i == artist.CreationDate && member == len(artist.Members) && j == FirstAlbumInt {
							Res = append(Res, artist)
						}
					}
				}
			}
		} else if Members == nil && Locations != "" {
			for i := CreationDateFrom; i <= CreationDateTo; i++ {
				for j := FirstAlbumFrom; j <= FirstAlbumTo; j++ {
					for Location := range artist.DatesLocations {
						FirstAlbumInt, _ := strconv.Atoi(artist.FirstAlbum[6:])
						if i == artist.CreationDate && Location == Locations && j == FirstAlbumInt {
							Res = append(Res, artist)
						}

					}
				}
			}
		} else {
			for i := CreationDateFrom; i <= CreationDateTo; i++ {
				for j := FirstAlbumFrom; j <= FirstAlbumTo; j++ {
					for Location := range artist.DatesLocations {
						for _, member := range MembersInt {
							FirstAlbumInt, _ := strconv.Atoi(artist.FirstAlbum[6:])
							if i == artist.CreationDate && Location == Locations && j == FirstAlbumInt && member == len(artist.Members) {
								Res = append(Res, artist)
							}
						}
					}
				}
			}
		}
	}
	tmpl, err := template.ParseFiles("./ui/html/home.html")
	if err != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	result := tmpl.Execute(w, Res)
	if result != nil {
		Errors(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	// fmt.Println(CreationDateFrom, CreationDateTo, FirstAlbumFrom, FirstAlbumTo, Members, Locations)
}

func Errors(w http.ResponseWriter, errorNum int, errorDescript string) {
	tmpl, err := template.ParseFiles("./ui/html/error.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errorNum)
	Error := Error{CodeError: errorNum, ErrorDescription: errorDescript}
	errors := tmpl.Execute(w, Error)
	if errors != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func checktxt(s string) bool {
	for _, v := range s {
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') || (v >= '0' && v <= '9') {
			return true
		}
	}
	return false
}

func unique(res []functions.Artist, i int) bool {
	for _, num := range res {
		if num.Id == i {
			return false
		}
	}
	return true
}
