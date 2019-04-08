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

	writer.Write([]string{"temp", "time", "site", "x", "y", "spin"})
	// output every k time frames
	nTimes := len(r[0].hist[0])
	k := int(math.Ceil(float64(nTimes-1) / float64(nOutTimes)))

	// loop over temperatures
	for _, scan := range r {
		T := scan.temp
		// loop over time frames of the first run
		for i := 0; i < len(scan.hist[0]); i += k {
			state := scan.hist[0][i]
			// loop over sites
			for id, site := range state {
				row := []string{
					strconv.FormatFloat(T, 'g', 5, 64),
					strconv.Itoa(i),
					strconv.Itoa(id),
					strconv.Itoa(site.Loc.X),
					strconv.Itoa(site.Loc.Y),
					strconv.Itoa(site.Spin)}
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
func exportMacroRecord(r []tempMacroHist, nOutTimes int) {
	filename := "macro_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"temp", "time", "mag", "ener"})
	// find the longest evolution time among scans
	var nTimes int
	for _, scan := range r {
		if len(scan.magHist) > nTimes {
			nTimes = len(scan.magHist)
		}
	}
	// output every k time frames
	k := math.Ceil(float64(nTimes-1) / float64(nOutTimes))

	// loop over scans
	for _, scan := range r {
		T := scan.temp
		// loop over time frames
		for i := 0; i < len(scan.magHist); i += int(k) {
			row := []string{
				strconv.FormatFloat(T, 'g', 5, 64),
				// unit of time is defined by number of sites
				strconv.Itoa(i),
				strconv.FormatFloat(scan.magHist[i], 'g', 5, 64),
				// unit of energy is defined by number of sites
				strconv.FormatFloat(scan.enerHist[i]/float64(mapDim*mapDim), 'g', 5, 64)}
			err := writer.Write(row)
			if err != nil {
				log.Fatal("Cannot write to file", err)
			}

		}
	}
	return
}
