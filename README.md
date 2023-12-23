# Batch MIDI Information

A tool to get quick MIDI information about a bunch of MIDI files, and save it to a JSON file.

## Purpose

This was created for getting the info from all of my MIDIs into a single file for my website. My website works off an object which follows a similar structure to the structure in the midis folder.

The program looks for all .mid files in the midis folder, and parses through all of them. It calculates:

- File Size
- Track Count
- PPQN
- Note Count

and saves it into a JSON file with the following structure:

```json
{
    "midiName": {
        "versionName": {
            "FileSize": string,
            "Tracks": int,
            "PPQN": int,
            "Notes": int
        },
        [...]
    }
    [...]
}
```

## Usage

To use this project, follow the steps below

1. Clone the repo
2. Install [Go](https://go.dev)
3. Put your MIDIs in the "midis" folder. **Note: [Use the structure prescribed in the README there](midis/).**
4. Execute `go run .` to run the program
5. The info will be spit out as "midi_info.json" in the root directory.

## Caveats

The code currently does not read note count accurately. I am working to solve this issue and get a more reliable note count but it may take me some time. If you are someone who has experience with note counts and want to help out, please make a pull request!
