package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/marni/goigc"
)

//StartTime ...
var StartTime time.Time

//MetaInf ..
var MetaInf []metaTrack

//Tracks ...
var Tracks []igc.Track

type metaInf struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

type trackIn struct {
	URL string `json:"url"`
}

type trackOut struct {
	ID int `json:"id"`
}

type metaTrack struct {
	Hdate       time.Time `json:"H_date"`
	Pilot       string    `json:"pilot"`
	GliderType  string    `json:"glider"`
	GliderID    string    `json:"glider_id"`
	TrackLength float64   `json:"track_length"`
}

func durationFormat(sec float64) string {
	var days, hours, minutes, seconds float64

	if sec > 86400 {
		seconds = math.Mod(sec, 86400.0)
		days = math.Trunc(sec / 86400.0)
		sec = seconds
	}
	if sec > 3600 {
		seconds = math.Mod(sec, 3600.0)
		hours = math.Trunc(sec / 3600.0)
		sec = seconds
	}
	if sec > 60 {
		seconds = math.Mod(sec, 60.0)
		minutes = math.Trunc(sec / 60.0)
		sec = seconds
	}
	upTime := "P" + strconv.FormatFloat(days, 'f', 0, 64) + "DT" + strconv.FormatFloat(hours, 'f', 0, 64) + "H" + strconv.FormatFloat(minutes, 'f', 0, 64) + "M" + strconv.FormatFloat(sec, 'f', 0, 64) + "S"
	return upTime
}

func handlerAPI(w http.ResponseWriter, r *http.Request) {
	d := time.Since(StartTime)
	dur := d.Seconds()
	upTime := durationFormat(dur)

	metadata := metaInf{upTime, "Service app for IGC tracks", "v1"}
	http.Header.Add(w.Header(), "content-type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}

func handlerTracksIn(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var input trackIn

	err := dec.Decode(&input)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	inurl, err := url.Parse(input.URL)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	igcsrc := inurl.String()
	tempTrack, err := igc.ParseLocation(igcsrc)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	Tracks = append(Tracks, tempTrack)
	id := len(Tracks)
	updateMeta(id)

	json.NewEncoder(w).Encode(id)
}

func handlerTracksOut(w http.ResponseWriter, r *http.Request) {
	id := len(Tracks)

	var idArray []int
	for i := 1; i <= id; i++ {
		idArray = append(idArray, i)
	}

	json.NewEncoder(w).Encode(idArray)
}

func handlerMetaTrack(w http.ResponseWriter, r *http.Request) {
	id, _ := findParams(r)
	if id >= len(Tracks) {
		http.Error(w, "Not Found", 404)
		return
	}

	json.NewEncoder(w).Encode(MetaInf[id])
}

func handlerSpecificTrack(w http.ResponseWriter, r *http.Request) {
	id, field := findParams(r)

	if id >= len(Tracks) {
		http.Error(w, "Not Found", 404)
		return
	}

	switch field {
	case "pilot":
		fmt.Fprintf(w, MetaInf[id-1].Pilot)
	case "glider":
		fmt.Fprintf(w, MetaInf[id-1].GliderType)
	case "glider_id":
		fmt.Fprintf(w, MetaInf[id-1].GliderID)
	case "track_length":
		fmt.Fprintf(w, strconv.FormatFloat(MetaInf[id-1].TrackLength, 'f', 0, 64))
	case "H_date":
		fmt.Fprintf(w, MetaInf[id-1].Hdate.String())
	default:
		http.Error(w, "Bad Request", 400)
		return
	}
}

func updateMeta(id int) {
	var tempMeta = metaTrack{
		Pilot:       Tracks[id-1].Header.Pilot,
		GliderType:  Tracks[id-1].Header.GliderType,
		GliderID:    Tracks[id-1].Header.GliderID,
		TrackLength: Tracks[id-1].Task.Distance(),
		Hdate:       Tracks[id-1].Header.Date,
	}
	MetaInf = append(MetaInf, tempMeta)
	fmt.Println(len(MetaInf), len(Tracks))
}

func findParams(r *http.Request) (int, string) {
	in := strings.Trim(r.URL.Path, "/lgcinfo/api/igc/")
	parts := strings.Split(in, "/")
	if len(parts) == 1 {
		id, err := strconv.Atoi(parts[0])
		if err != nil {
			//error
			return 0, ""
		}
		return id, ""
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		//error
		return 0, ""
	}
	field := parts[1]
	return id, field
}

//GetPort ...
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "5000"
		fmt.Println("Could not find port in environment, setting port to: " + port)
	}
	return ":" + port
}

func main() {
	StartTime = time.Now()
	r := mux.NewRouter()
	r.HandleFunc("/igcinfo/api", handlerAPI).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc", handlerTracksIn).Methods("POST")
	r.HandleFunc("/igcinfo/api/igc", handlerTracksOut).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc/{id}", handlerMetaTrack).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc/{id}/{field}", handlerSpecificTrack).Methods("GET")
	http.ListenAndServe(GetPort(), r)
}
