package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

// write state values of the player to a csv file
func exportMagHist(hist []Situation) {
	filename := "mag_hist.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for time, situ := range hist {
		row := []string{
			strconv.Itoa(time),
			situ.action,
			strconv.FormatFloat(situ.mag, 'g', 5, 64)}
		err := writer.Write(row)
		if err != nil {
			log.Fatal("Cannot write to file", err)
		}
	}

	return
}
