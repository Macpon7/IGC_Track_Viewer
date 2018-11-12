package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/marni/goigc"
)

//StartTime ...
var StartTime time.Time

//MetaInf ..
var MetaInf []metaTrack

//totalID ...
var totalID int

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

	totalID++

	tempLength := 0.0
	for i := 0; i < len(tempTrack.Points)-1; i++ {
		tempLength += tempTrack.Points[i].Distance(tempTrack.Points[i+1])
	}

	var tempMeta = metaTrack{
		Pilot:       tempTrack.Header.Pilot,
		GliderType:  tempTrack.Header.GliderType,
		GliderID:    tempTrack.Header.GliderID,
		TrackLength: tempLength,
		Hdate:       tempTrack.Header.Date,
	}
	MetaInf = append(MetaInf, tempMeta)

	json.NewEncoder(w).Encode(totalID)
}

func handlerTracksOut(w http.ResponseWriter, r *http.Request) {
	var idArray []int
	for i := 1; i <= totalID; i++ {
		idArray = append(idArray, i)
	}

	json.NewEncoder(w).Encode(idArray)
}

func handlerMetaTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varID := vars["id"]
	id, err := strconv.Atoi(varID)
	if err != nil {
		//error
		return
	}

	if id > totalID {
		http.Error(w, "Not Found", 404)
		return
	}

	json.NewEncoder(w).Encode(MetaInf[id-1])
}

func handlerSpecificTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	varID := vars["id"]
	field := vars["field"]
	id, err := strconv.Atoi(varID)
	if err != nil {
		//error
		return
	}

	if id > totalID {
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

//GetPort ...
func GetPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("Could not find port in environment, setting port to: " + port)
	}
	return ":" + port
}

func init() {
	totalID = 0
	StartTime = time.Now()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/igcinfo/api", handlerAPI).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc", handlerTracksIn).Methods("POST")
	r.HandleFunc("/igcinfo/api/igc", handlerTracksOut).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc/{id}", handlerMetaTrack).Methods("GET")
	r.HandleFunc("/igcinfo/api/igc/{id}/{field}", handlerSpecificTrack).Methods("GET")
	http.ListenAndServe(GetPort(), r)
}
