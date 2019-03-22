package main

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"

	pb "gopkg.in/cheggaaa/pb.v1"
)

// print out at most nOutTimes equally spread time frames along state history
func exportStateHist(h []State, nOutTimes int) {
	filename := "state_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"time", "site", "x", "y", "conns", "spin"})
	nSites := len(h[0].Spins)
	nTimes := len(h)
	k := math.Ceil(float64(nTimes) / float64(nOutTimes)) // output every k time frames

	bar := pb.StartNew(nTimes)
	for i, state := range h { // loop over times
		bar.Increment()
		if math.Mod(float64(i), k) == 0.0 {
			for id, loc := range state.Locations { // loop over sites
				row := []string{
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
	bar.FinishPrint("state history exported")
	return
}

// print out at most nOutTimes equally spread time frames along magnetization history
func exportMagHist(h []float64, nSites, nOutTimes int) {
	filename := "mag_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"time", "mag"})
	nTimes := len(h)
	k := math.Ceil(float64(nTimes) / float64(nOutTimes)) // output every k time frames

	bar := pb.StartNew(nTimes)
	for i, mag := range h {
		bar.Increment()
		if math.Mod(float64(i), k) == 0.0 {
			row := []string{
				// unit of time is defined by number of sites
				strconv.FormatFloat(float64(i)/float64(nSites), 'g', 5, 64),
				strconv.FormatFloat(mag, 'g', 5, 64)}
			err := writer.Write(row)
			if err != nil {
				log.Fatal("Cannot write to file", err)
			}
		}
	}
	bar.FinishPrint("magnetization history exported")
	return
}
