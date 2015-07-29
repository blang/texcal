package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"
	"time"
)

var DAY time.Duration = time.Second * 86400

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	os.Exit(1)
}

var monthStrM = map[string]map[time.Month]string{
	"de": {
		time.January:   "Januar",
		time.February:  "Februar",
		time.March:     "März",
		time.April:     "April",
		time.May:       "Mai",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "August",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Dezember",
	},
}

var startDateStr = flag.String("begin", "", "Date to start (02.01.2006)")
var endDateStr = flag.String("end", "", "Date to end (02.01.2006)")
var argDays = flag.Int("days", 30, "Number of days")
var outputFileStr = flag.String("output", "", "Outputfile, otherwise print to stdout")
var langStr = flag.String("lang", "en", "Language")

var fmap = template.FuncMap{
	"mod": func(i int, j int) int {
		return i % j
	},
	"day": func(t time.Time) int {
		return t.Day()
	},

	"month": func(t time.Time) string {
		if *langStr != "en" {
			return monthStrM[*langStr][t.Month()]
		}
		return t.Month().String()
	},
}

func main() {
	flag.Parse()

	// Startdate
	if *startDateStr == "" {
		*startDateStr = time.Now().Format("02.01.2006")
	}
	startDate, err := time.Parse("02.01.2006", *startDateStr)
	if err != nil {
		fatal("Error parsing start date: %s\n", err)
	}

	// Enddate
	if *endDateStr == "" {
		*endDateStr = startDate.Add(time.Duration(*argDays) * DAY).Format("02.01.2006")
	}
	endDate, err := time.Parse("02.01.2006", *endDateStr)
	if err != nil {
		fatal("Error parsing end date: %s\n", err)
	}
	if startDate.After(endDate) {
		fatal("Startdate after enddate!")
	}

	tmpl, err := template.New("").Funcs(fmap).ParseFiles("calendar.tex")
	if err != nil {
		fatal("Error parsing template: %s\n", err)
	}

	// Check if startDate is Monday, otherwise find previous Monday
	for startDate.Weekday() != time.Monday {
		startDate = startDate.Add(-DAY)
	}

	// Check if endDate is Sunday, otherwise find next Sunday
	for endDate.Weekday() != time.Sunday {
		endDate = endDate.Add(DAY)
	}

	// Build slice of dates to print
	var dates []time.Time
	for startDate.Before(endDate) {
		dates = append(dates, startDate)
		startDate = startDate.Add(DAY)
	}
	dates = append(dates, endDate)

	var writer io.Writer
	writer = os.Stdout // Default to stdout
	if *outputFileStr != "" {
		_, err := os.Stat(*outputFileStr)
		if err == nil {
			fatal("Output file already exists: %s\n", *outputFileStr)
		}
		fout, err := os.OpenFile(*outputFileStr, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			fatal("Could not open output file: %s\n", err)
		}
		defer fout.Close()
		writer = fout
	}

	err = tmpl.ExecuteTemplate(writer, "calendar.tex", dates)
	if err != nil {
		fatal("Error executing template: %s\n", err)
	}
}
