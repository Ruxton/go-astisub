package main

import (
	stlflag "flag"
	"io"
	"log"
	"os"

	"github.com/asticode/go-subtitles"
	"github.com/asticode/go-toolkit/flag"
)

// Flags
var (
	duration   = stlflag.Duration("d", 0, "the input duration")
	inputPath  = stlflag.String("i", "", "the input path")
	outputPath = stlflag.String("o", "", "the output path (optionnal)")
)

func main() {
	// Parse flags
	var subcommand = flag.Subcommand()
	stlflag.Parse()

	// Switch on subcommand
	switch subcommand {
	case "add":
		// Open file
		var s *subtitles.Subtitles
		var err error
		if s, err = subtitles.Open(*inputPath); err != nil {
			log.Fatal(err)
		}

		// Add duration
		s.Add(*duration)

		// Init the writer
		var w io.Writer
		if *outputPath != "" {
			// Create output file
			var file *os.File
			if file, err = os.Create(*outputPath); err != nil {
				log.Fatal(err)
			}
			w = file
		} else {
			w = os.Stdout
		}

		// Output the content
		if err = subtitles.ToWriterSRT(*s, w); err != nil {
			log.Fatal(err)
		}
	}
}
