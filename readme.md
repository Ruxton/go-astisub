This is a Golang library to manipulate subtitles files.

The CLI let you add a positive or negative duration to each timestamps of your subtitles file, allowing you to fix delays/hurries in your subtitle file.

# Installation

Run the following command :

    go get github.com/asticode/go-subtitles && go install github.com/asticode/go-subtitles/cmd/main.go
    
This will fetch the appropriate dependencies and install the CLI version of `go-subtitles`

# Use the CLI

Provided that you setup your Go environment correctly, you now have a basic `go-subtitles` command to interact with subtitle files.

    Usage:
      -d duration
        	the input duration
      -i string
        	the input path
      -o string
        	the output path (optionnal)

Hence to add `123ms` to the input path `/path/to/input` writing to the output path `/path/to/output` :

    go-subtitles -i /path/to/input -d 123ms -o /path/to/output
