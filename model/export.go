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
	bar := pb.StartNew(len(h))
	for time, state := range h { // loop over times
		bar.Increment()
		if math.Mod(float64(time), float64(k)) == 0.0 {
			for id, loc := range state.Locations { // loop over sites
				row := []string{
					strconv.Itoa(time),
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

// write situation history every k time frames
func exportSituHist(h []Situation, k int) {
	filename := "situ_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"time", "action", "magnetization"})
	bar := pb.StartNew(len(h))
	for time, situ := range h {
		bar.Increment()
		if math.Mod(float64(time), float64(k)) == 0.0 {
			row := []string{
				strconv.Itoa(time),
				situ.Action,
				strconv.FormatFloat(situ.Mag, 'g', 5, 64)}
			err := writer.Write(row)
			if err != nil {
				log.Fatal("Cannot write to file", err)
			}
		}
	}
	bar.FinishPrint("situation history exported")
	return
}
