# LostPetsApp API

This project was Developed in GO.

Current rountes include
```json
{
    "method": "GET",
    "path": "/sightings/private/:guid",
    "name": "lostpets/internal/http.(*sightingsHandler).handleGetByGUID.func1"
  },
  {
    "method": "GET",
    "path": "/",
    "name": "lostpets/internal/http.apiInfoHandler.func1"
  },
  {
    "method": "GET",
    "path": "/sightings/:id",
    "name": "lostpets/internal/http.(*sightingsHandler).handleGetByID.func1"
  },
  {
    "method": "GET",
    "path": "/sightings/private/:guid/matches",
    "name": "lostpets/internal/http.(*sightingsHandler).handleGetAllMatches.func1"
  },
  {
    "method": "GET",
    "path": "/pet-types",
    "name": "lostpets/internal/http.getPetTypesHandler.func1"
  },
  {
    "method": "GET",
    "path": "/postings/:id",
    "name": "lostpets/internal/http.(*postingsHandler).handleGetByID.func1"
  },
  {
    "method": "GET",
    "path": "/postings/private/:guid",
    "name": "lostpets/internal/http.(*postingsHandler).handleGetByGUID.func1"
  },
  {
    "method": "GET",
    "path": "/postings",
    "name": "lostpets/internal/http.(*postingsHandler).handleGetAll.func1"
  },
  {
    "method": "POST",
    "path": "/postings",
    "name": "lostpets/internal/http.(*postingsHandler).handleCreatePosting.func1"
  },
  {
    "method": "GET",
    "path": "/pet-pictures/:id",
    "name": "lostpets/internal/http.(*fileHandler).handleServeFile.func1",
    "note": "NOT TESTED"
  },
  {
    "method": "POST",
    "path": "/pet-pictures",
    "name": "lostpets/internal/http.(*fileHandler).handleUploadFile.func1",
    "note": "NOT TESTED"
  },
  {
    "method": "GET",
    "path": "/postings/private/:guid/matches",
    "name": "lostpets/internal/http.(*postingsHandler).handleGetAllMatches.func1"
  },
  {
    "method": "GET",
    "path": "/sightings",
    "name": "lostpets/internal/http.(*sightingsHandler).handleGetAll.func1"
  },
  {
    "method": "POST",
    "path": "/sightings",
    "name": "lostpets/internal/http.(*sightingsHandler).handleCreateSighting.func1"
  }
]
```

the initial setup for file upload has been added, just not tested

make sure a db matching the name set int he config is created before running migrations

can use make file to build migration manager and to build the api

to run the api and migration manager, need the '-c' flag followed by the path to the config