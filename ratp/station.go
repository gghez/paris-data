package ratp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// StopCoordinates gathers latitude and longitude data
type StopCoordinates [2]float64

// Lat stands for Latitude
func (c StopCoordinates) Lat() float64 {
	return c[0]
}

// Lng stands for Longitude
func (c StopCoordinates) Lng() float64 {
	return c[1]
}

// StopRecord identifies a station geolocalized data
type StopRecord struct {
	Coordinates StopCoordinates `json:"stop_coordinates"`
	StopName    string          `json:"stop_name"`
	ID          string          `json:"stop_id"`
}

type dataSet struct {
	Nhits   int
	Records []struct {
		Fields StopRecord
	}
}

func (d dataSet) getStopRecords() (records []StopRecord) {
	for _, r := range d.Records {
		records = append(records, r.Fields)
	}
	return
}

const apiGrabLimit int = 10000
const requestPageSize int = 1000

func getRecordsPage(startIndex, pageSize int, chRecords chan<- []StopRecord, chCount chan<- int) {
	url := fmt.Sprintf("https://data.ratp.fr/api/records/1.0/search/?dataset=positions-geographiques-des-stations-du-reseau-ratp&facet=stop_name&rows=%d&start=%d", pageSize, startIndex)
	log.Printf("GET %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var ds dataSet
	json.NewDecoder(resp.Body).Decode(&ds)

	if chCount != nil {
		log.Printf("Total elements: %d\n", ds.Nhits)
		chCount <- ds.Nhits
	}

	chRecords <- ds.getStopRecords()
}

// GetMetroStops retrieves RATP metro stations data
func GetMetroStops() []StopRecord {
	chCount := make(chan int)
	chRecords := make(chan []StopRecord)
	go getRecordsPage(0, requestPageSize, chRecords, chCount)
	total := <-chCount

	maxRecordCount := apiGrabLimit
	if total < apiGrabLimit {
		maxRecordCount = total
	}

	startIndex, pageSize := requestPageSize, requestPageSize
	for startIndex < maxRecordCount {
		if startIndex+pageSize > maxRecordCount {
			pageSize = maxRecordCount - startIndex
		}

		go getRecordsPage(startIndex, pageSize, chRecords, nil)
		startIndex += pageSize
	}

	records := make([]StopRecord, 0, maxRecordCount)
	for len(records) < maxRecordCount {
		records = append(records, <-chRecords...)
	}

	return records
}
