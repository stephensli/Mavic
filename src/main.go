package main

import (
	"log"
	"os"
	"strings"

	"github.com/tehstun/mavic/src/reddit"
	"github.com/urfave/cli"
)

var app = cli.NewApp()
var options = reddit.Options{}

func setupApplicationInformation() {
	app.Name = "Mavic"
	app.Description = "Mavic is a CLI application designed to download direct images found on selected reddit subreddits."
	app.Usage = ".\\mavic.exe --subreddits cute -l 100 --output ./pictures -f"
	app.Author = "Stephen Lineker-Miller <slinekermiller@gmail.com>"
	app.Version = "0.0.1"
}

func setupApplicationFlags() {
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "output, o",
			Usage:       "The output directory to store the images.",
			Value:       "./",
			Destination: &options.OutputDirectory,
		}, cli.IntFlag{
			Name:        "limit, l",
			Usage:       "The total number of posts max per sub-reddit",
			Value:       50,
			Destination: &options.ImageLimit,
		},
		cli.BoolFlag{
			Name:        "frontpage, f",
			Usage:       "If the front page should be scrapped or not.",
			Destination: &options.FrontPage,
		},
		cli.StringFlag{
			Name:        "type, t",
			Usage:       "What kind of page type should reddit be during the scrapping process. e.g hot, new. top.",
			Value:       "hot",
			Destination: &options.PageType,
		},
		cli.StringSliceFlag{
			Name:     "subreddits, s",
			Usage:    "What subreddits are going to be scrapped for downloading images.",
			Required: true,
		},
	}
}

// processSubreddits takes in a ring of possible sub reddits and splits them into
// a slice of the sub reddits to be processed, there is currently a bug with the
// cli tools which is resulting in the funky processing and its best to just
// process it as a string for the time being.
func processSubreddits(subreddits string, arguments []string) []string {
	// since it only seems to parse the first element, even though more was selected
	// so we push it here and then go and grab the remaining.
	processed := []string{subreddits}

	for i := 0; i < len(arguments); i++ {
		value := arguments[i]

		// if we have hit the next command, then we must breakout since we
		// no longer have any more subs.
		if strings.HasPrefix(value, "-") {
			break
		}

		processed = append(processed, value)
	}

	return processed
}

// start is called by the cli control when the cli controls are parsed, setting up
// and building a context around the cli application. This is the time the sub
// reddits are parsed since the cli tools don't support binding stringSlices.
func start(c *cli.Context) error {
	options.Subreddits = processSubreddits(c.String("subreddits"), c.Args())

	// if it equals nil, and no sub reddits was given, then just set them
	// as s empty slice, letting the scraper handle the empty case as
	// it should.
	if options.Subreddits == nil {
		options.Subreddits = []string{}
	}

	// create a new reddit scraper and process through all the sub reddits
	// downloading the images in the output folder / sub reddit / image.
	scraper := reddit.NewScraper(options)
	scraper.ProcessSubreddits()
	return nil
}

func main() {
	setupApplicationInformation()
	setupApplicationFlags()

	app.Action = start
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
