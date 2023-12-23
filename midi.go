package main

import (
	"errors"
	"fmt"
	"os"
)

func parseMIDIVersion(midi string) (Stat, error) {
	logf("parsing midi: %v", midi)
	stat := Stat{}
	// open midi file
	midiFile, err := os.Open(midi)
	if err != nil {
		return stat, err
	}
	defer midiFile.Close()

	// get file size
	fileInfo, err := midiFile.Stat()
	if err != nil {
		return stat, err
	}
	stat.FileSize = ByteCountSI(fileInfo.Size())

	// parse header track

	// parse type
	// ensure that type is of MThd
	headerType := make([]byte, 4)
	_, err = midiFile.Read(headerType)
	if err != nil {
		return stat, err
	}

	if string(headerType) != "MThd" {
		logf("invalid header track | header type: %v", string(headerType))
		return stat, errors.New("MIDI file does not contain header track")
	}

	// parse header size
	// ensure that header size is 6
	headerSize := make([]byte, 4)
	_, err = midiFile.Read(headerSize)
	if err != nil {
		return stat, err
	}
	headerSizeInt := int(headerSize[0])<<24 | int(headerSize[1])<<16 | int(headerSize[2])<<8 | int(headerSize[3])

	if headerSizeInt != 6 {
		logf("invalid header size (>6) | header type: %v", string(headerType))
		return stat, errors.New("MIDI header size is not 6")
	}

	// parse format
	// ensure that format is 1
	format := make([]byte, 2)
	_, err = midiFile.Read(format)
	if err != nil {
		return stat, err
	}

	if format[0] != 0 || format[1] != 1 {
		logf("invalid midi format | header type: %v", string(headerType))
		return stat, errors.New("MIDI format is not 1")
	}

	// parse track count
	// write to trackCountInt as an integer
	trackCount := make([]byte, 2)
	_, err = midiFile.Read(trackCount)
	if err != nil {
		return stat, err
	}

	trackCountInt := int(trackCount[0])<<8 | int(trackCount[1])

	// parse time division
	// convert to int
	timeDivision := make([]byte, 2)
	_, err = midiFile.Read(timeDivision)
	if err != nil {
		return stat, err
	}

	timeDivisionInt := int(timeDivision[0])<<8 | int(timeDivision[1])

	// save header info
	logf("parsed header: %v (%v bytes) with %v tracks and time division of %v", string(headerType), headerSizeInt, trackCountInt, timeDivisionInt)
	stat.Tracks = trackCountInt
	stat.PPQN = timeDivisionInt

	// parse track chunks for notes
	logf("parsing tracks for notes")
	for i := 0; i < trackCountInt; i++ {
		// parse track header
		// ensure that track header is MTrk
		trackType := make([]byte, 4)
		_, err = midiFile.Read(trackType)
		if err != nil {
			return stat, err
		}

		if string(trackType) != "MTrk" {
			logf("invalid header track | header type: %v", string(trackType))
			return stat, errors.New("MIDI file does not contain header track")
		}

		// parse track size
		// convert to int
		trackSize := make([]byte, 4)
		_, err = midiFile.Read(trackSize)
		if err != nil {
			return stat, err
		}

		trackSizeInt := int(trackSize[0])<<24 | int(trackSize[1])<<16 | int(trackSize[2])<<8 | int(trackSize[3])

		// while there is track data to read
		// read event type
		// if event type is note on
		// add to note count
		trackData := make([]byte, trackSizeInt)
		n, err := midiFile.Read(trackData)
		if err != nil {
			return stat, err
		}
		offset := 0
		loop := true
		for offset < n && loop {
			dt := 0
			// read delta time
			parseDeltaTime(&trackData, &offset, &dt)

			if offset >= n {
				break
			}

			// read command
			command := trackData[offset]
			offset++

			switch command & 0xF0 {
			case 0x90:
				{
					// Note on
					stat.Notes++
					// Skip 2 bytes
					offset += 2
					break
				}
			case 0xA0:
			case 0xB0:
			case 0xE0:
			case 0x80:
				{
					// Skip 2 bytes
					offset += 2
					break
				}
			case 0xC0:
			case 0xD0:
				{
					// Skip 1 byte
					offset += 1
					break
				}
			default:
				{
					switch command {
					case 0xFF:
						{
							metaEvent := trackData[offset]
							offset++
							switch metaEvent {
							case 0x00:
							case 0x59:
								// FF nn 02 ss ss
								// Sequence Number and Key Signature
								// Skip all of this
								offset += 4
							case 0x01:
							case 0x02:
							case 0x03:
							case 0x04:
							case 0x05:
							case 0x06:
							case 0x07:
							case 0x7F:
								// FF 0n ll tt tt tt tt ...
								// Text Event
								// Skip all of this
								dt = 0
								parseDeltaTime(&trackData, &offset, &dt)
								offset += dt
							case 0x20:
								// FF 20 01 cc
								// MIDI Channel Prefix
								// Skip all of this
								offset += 3
							case 0x2F:
								// FF 2F 00
								// End of Track
								// Skip all of this
								offset += 2
								loop = false
							case 0x51:
								// FF 51 03 tt tt tt
								// Set Tempo
								// Skip all of this
								offset += 5
							case 0x54:
								// FF 54 05 hr mn se fr ff
								// SMPTE Offset
								// Skip all of this
								offset += 7
							case 0x58:
								// FF 58 04 nn dd cc bb
								// Time Signature
								// Skip all of this
								offset += 6
							}
						}
					case 0xF0:
					case 0xF7:
						{
							// Skip variable length data
							dt = 0
							parseDeltaTime(&trackData, &offset, &dt)
							offset += dt
							break
						}
					case 0xF2:
						{
							// Skip 2 bytes
							offset += 2
						}
					case 0xF3:
						{
							// Skip 1 byte
							offset += 1
						}
					}
				}
			}
		}
	}
	logf("parsed tracks for notes. found %v notes", stat.Notes)

	return stat, nil
}

func parseDeltaTime(data *[]byte, offset *int, dt *int) {
	*dt = 0
	for {
		val := int((*data)[*offset])
		*offset++
		*dt = (*dt << 7) | (val & 0x7F)
		if val&0x80 == 0 {
			break
		}
	}
}

func logf(format string, a ...any) {
	println(fmt.Sprintf(format, a...))
}
