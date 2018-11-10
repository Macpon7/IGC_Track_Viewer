package main

import (
	"net/http"
)

func handlerMetaAPI(w http.ResponseWriter, r *http.Request) {
	/*
		~~GET~~
		Return meta information about the API
			"uptime": <uptime>, (formatted to ISO 8601)
			"info": "Service for IGC tracks.",
			"version": "v1",
		response-type: application/json
		response code: 200
	*/
}

func handlerTracks(w http.ResponseWriter, r *http.Request) {
	/*
		~~POST~~
		register tracks, response-type: application/json, "handle all errors gracefully" ;)
		request template:
			"url": "<url>"  <- url is a normal url
		response template:
			"id": "<id>" <- unique id to that specific file upload

		~~GET~~
		return the array of all track ID's
			[<id1>, <id2>, ...]
		response-type: application/json
		response code: 200, or appropriate error code
		response: the array of ID's, or an empty array if no tracks have been stored yet
	*/
}

func handlerMetaTrack(w http.ResponseWriter, r *http.Request) {
	/*
		~~GET~~ /api/igc/<id>
		return the meta information about a track with provided <id>, or NOT FOUND code
		application/json
		200 or appropriate
		Response:
			"H_date": 		<date from File Header, H-record>,
			"pilot":		<pilot>,
			"glider":		<glider>,
			"glider_id":	<glider_id>,
			"track_length":	<calculated total track length>

		~~GET~~ /api/igc/<id>/<field>
		return the single detailed meta information or NOT FOUND
		response should always be a string, with the exception of the calculated total length
	*/
}

func main() {

}
