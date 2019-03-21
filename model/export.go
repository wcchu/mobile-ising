package main

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"

	pb "gopkg.in/cheggaaa/pb.v1"
)

// write state history every k time frames
func exportStateHist(h []State, k int) {
	filename := "state_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"time", "site", "x", "y", "conns", "spin"})
	n := len(h[0].Spins)
	bar := pb.StartNew(len(h))
	for i, state := range h { // loop over times
		bar.Increment()
		if math.Mod(float64(i), float64(k)) == 0.0 {
			for id, loc := range state.Locations { // loop over sites
				row := []string{
					// time is normalized by number of sites
					strconv.FormatFloat(float64(i)/float64(n), 'g', 5, 64),
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

// write magnetization history every k time frames
func exportMagHist(h []float64, n, k int) {
	filename := "mag_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"time", "mag"})
	bar := pb.StartNew(len(h))
	for i, mag := range h {
		bar.Increment()
		if math.Mod(float64(i), float64(k)) == 0.0 {
			row := []string{
				// time is normalized by number of sites
				strconv.FormatFloat(float64(i)/float64(n), 'g', 5, 64),
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
