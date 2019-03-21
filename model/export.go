package main

import (
	"encoding/csv"
	"log"
	"math"
	"os"
	"strconv"
)

// write situation history
func exportMagHist(hist []Situation) {
	filename := "situ_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for time, situ := range hist {
		if math.Mod(float64(time), 10.0) == 0.0 {
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
	return
}
