package main

import (
	"encoding/json"
	"os"
	"strings"
)

func main() {
	println("Batch MIDI Information")
	println("Created by: 6gh")
	println("Report bugs to: https://github.com/6gh/Batch-MIDI-Information/issues")
	println("")

	// check if midis folder exists
	if _, err := os.Stat("midis"); os.IsNotExist(err) {
		println("Error: midis folder not found")
		return
	}

	// find midis
	midis, err := findMIDIs("midis")
	if err != nil {
		println("Error:", err.Error())
		return
	}

	midiNames := map[string]map[string]Stat{}
	for _, midi := range midis {
		// split path by / if on linux or \ if on windows
		split := strings.Split(midi, "/")
		if len(split) == 1 {
			split = strings.Split(midi, "\\")
		}

		// get midi name
		midiName := split[len(split)-2]

		// remove .mid
		versionName := strings.Replace(split[len(split)-1], ".mid", "", 1)

		// parse midi
		stat, err := parseMIDIVersion(midi)
		if err != nil {
			println("Error:", err.Error())
			return
		}

		// check if midi name is already in map, if not add it
		if _, ok := midiNames[midiName]; !ok {
			midiNames[midiName] = map[string]Stat{}
		}

		// add midi to map
		midiNames[midiName][versionName] = stat
	}

	// save midi info as json
	j, err := json.MarshalIndent(midiNames, "", "  ")
	if err != nil {
		println("Error:", err.Error())
		return
	}

	// print midi info
	println(string(j))

	err = os.WriteFile("midi_info.json", j, 0644)
	if err != nil {
		println("Error:", err.Error())
		return
	}

	println("Saved midi info to midi_info.json")
}
