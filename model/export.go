package main

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

// print out at most nOutTimes equally spread time frames along state history
func exportStateRecord(r []tempStateHist, nOutTimes int) {
	filename := "state_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"temp", "time", "site", "x", "y", "conns", "spin"})
	nTimes := len(r[0].hist) // total time defined by the lowest-temp scan
	nSites := len(r[0].hist[0].Spins)
	k := math.Ceil(float64(nTimes) / float64(nOutTimes)) // output every k time frames

	for _, scan := range r { // loop over scans
		T := scan.temp
		for i := 0; i < len(scan.hist); i += int(k) { // loop over times
			state := scan.hist[i]
			for id, loc := range state.Locations { // loop over sites
				row := []string{
					strconv.FormatFloat(T, 'g', 5, 64),
					// unit of time is defined by number of sites
					strconv.FormatFloat(float64(i)/float64(nSites), 'g', 5, 64),
					strconv.Itoa(id),
					strconv.FormatFloat(loc.X, 'g', 5, 64),
					strconv.FormatFloat(loc.Y, 'g', 5, 64),
					strconv.Itoa(state.Connections[id]),
					strconv.Itoa(state.Spins[id])}
				err := writer.Write(row)
				if err != nil {
					log.Fatal("Cannot write to file", err)
				}
			}
		}
	}

	return
}

// print out at most nOutTimes equally spread time frames along magnetization history
func exportMagRecord(r []tempMagHist, nSites, nOutTimes int) {
	filename := "mag_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"temp", "time", "mag"})
	nTimes := len(r[0].hist)
	k := math.Ceil(float64(nTimes) / float64(nOutTimes)) // output every k time frames

	for _, scan := range r {
		T := scan.temp
		for i := 0; i < len(scan.hist); i += int(k) {
			row := []string{
				strconv.FormatFloat(T, 'g', 5, 64),
				// unit of time is defined by number of sites
				strconv.FormatFloat(float64(i)/float64(nSites), 'g', 5, 64),
				strconv.FormatFloat(scan.hist[i], 'g', 5, 64)}
			err := writer.Write(row)
			if err != nil {
				log.Fatal("Cannot write to file", err)
			}

		}
	}
	return
}
