# IGC_Track_Viewer
An in-memory web app for viewing igc files. This was Assignment 1 for IMT2681 at NTNU Gjøvik during the fall of 2018.
The application can be accesed at http://igctrackview.herokuapp.com/

## Features
Able to import igc files, store relevant information about the igc files in memory, and return to the user that information through API calls.


## API Reference
These are the api calls that can be made to http://igctrackview.herokuapp.com/igcinfo to upload an igc file or get some information from a previously uploaded file

```
GET /api
```
Returns meta information about the API.
```
POST /api/igc
```
Request body must contain a link to an igc file to properly upload it. Once the file has been parsed the unique ID of the file is returned, which can be used to access information about that specific file.
```
GET /api/igc
```
Returns an array containing the ID of all tracks currently stored in memory.
```
GET /api/igc/<id>
```
Returns meta information about the track with the given id.
```
GET /api/igc/<id>/<field>
```
Returns a specific field of information about the track with the given id. The field can be: "pilot", "glider", "glider_id", "H_date", or "track_length".

## Credits
IGC library for the processing of igc files: [goigc](https://github.com/marni/goigc)