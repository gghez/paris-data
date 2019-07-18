/*

This script get RATP (National Transport Company in France) stations stop geo information

Opendata API URL is: https://data.ratp.fr/api/records/1.0/search/?dataset=positions-geographiques-des-stations-du-reseau-ratp&facet=stop_name

The record JSON structure in response payload is:
{
	"datasetid": "positions-geographiques-des-stations-du-reseau-ratp",
	"recordid": "2961365fdfdc35efa97ce4bc88f8b1a64c4adc8a",
	"fields": {
		"stop_coordinates": [
			48.840882627372274,
			2.549350971449035
		],
		"stop_desc": "Esplanade des Arcades - 93051",
		"stop_name": "Noisy-le-Grand (Mont d'Est)",
		"stop_id": "1652"
	},
	"geometry": {
		"type": "Point",
		"coordinates": [
			2.549350971449035,
			48.840882627372274
		]
	},
	"record_timestamp": "2019-03-22T10:43:58+00:00"
}

*/

package main

import (
	"fmt"

	ratp "github.com/gghez/paris-data/ratp"
)

func main() {
	fmt.Printf("Grabbed %d records.\n", len(ratp.GetMetroStops()))
}
