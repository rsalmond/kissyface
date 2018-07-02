package cmd

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"strings"
	"time"
)

func parseArgs() (filename string, err error) {
	if len(os.Args) != 2 {
		return "", errors.New(fmt.Sprintf("Usage: %s \"<filename>\"\n", os.Args[0]))
	} else {
		return os.Args[1], nil
	}
}

type Message struct {
	Body   string
	User   string
	Day    int
	Month  int
	Year   int
	Hour   int
	Minute int
	Second int
}

type Histogram struct {
	Hours         []int
	Hourly        map[time.Time]map[string]int
	HourlyOrder   []time.Time
	Weekdays      map[time.Weekday]int
	Users         map[string]int
	TotalMessages int
}

func (h *Histogram) init() {
	// prepare all the histogram data
	h.Hours = make([]int, 24)
	h.Hourly = make(map[time.Time]map[string]int)
	h.HourlyOrder = make([]time.Time, 0)
	h.Weekdays = make(map[time.Weekday]int, 7)
	h.Users = make(map[string]int, 2)
	h.TotalMessages = 0

	for i := 0; i < 24; i++ {
		h.Hours[i] = 0
	}

	// june 17th was a sunday, this sets up our slice in weekday order
	for date := time.Date(2018, time.Month(6), 17, 12, 0, 0, 0, time.UTC); date.Day() < 25; date = date.Add(time.Hour * time.Duration(24)) {
		h.Weekdays[date.Weekday()] = 0
	}
}

func (h *Histogram) count(m *Message) {
	// count every mesage
	h.TotalMessages++
	// sum up total messages by hour of the day they were sent
	h.Hours[m.Hour]++

	// and by day of the week
	date := time.Date(m.Year, time.Month(m.Month), m.Day, m.Hour, 0, 0, 0, time.UTC)
	h.Weekdays[date.Weekday()]++

	// and by messages per hour across all time (broken out by user)
	if _, date_present := h.Hourly[date]; date_present {
		if _, user_present := h.Hourly[date][m.User]; user_present {
			h.Hourly[date][m.User]++
		} else {
			h.Hourly[date][m.User] = 1
		}
	} else {
		h.Hourly[date] = make(map[string]int, 2)
		h.Hourly[date][m.User] = 1
	}
	// store these hourly increments in an ordered slice so we can write out the map in order later
	h.HourlyOrder = append(h.HourlyOrder, date)

	// and by user who sent them
	if _, present := h.Users[m.User]; present {
		h.Users[m.User]++
	} else {

		h.Users[m.User] = 1
	}

	return
}

func (h Histogram) write_alltime_csv() (error) {
	const filename string = "./all_time_by_hour.csv"

	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write csv file %s", filename))
	}

	defer f.Close()

	// this will hold the mapping for each user to the number of messages for the current row of output
	current_row := make(map[string]int)

	// used to populate the headers and individual data rows
	usernames := make([]string, 0)
	var messages []string

	// initialize
	for user, _ := range h.Users {
		current_row[user] = 0
		usernames = append(usernames, user)
	}

	// write header
	f.WriteString(fmt.Sprintf("Hour, %s\n", strings.Join(usernames, ",")))

	// get the first and last elements from the ordered slice of all messages for start / end times
	// probably don't even need this whole slice but whatever
	start_time := h.HourlyOrder[0]
	end_time := h.HourlyOrder[len(h.HourlyOrder)-1]

	// loop over every hour from the time of the very first message
	for current_hour := start_time; current_hour.Before(end_time); current_hour = current_hour.Add(time.Hour * 1) {
		// reset the message slice for each row we write
		messages = make([]string, 0)
		// check if we have any messages counted in this hour
		if _, ok := h.Hourly[current_hour]; ok {
			// if we do, record the count for that user in the current row
			for user, messages := range h.Hourly[current_hour] {
				current_row[user] = messages
			}

		}
		// generate the CSV string values for this row of data (and reset the counts in preparation for the next row)
		for _, user := range usernames {
			messages = append(messages, fmt.Sprintf("%d", current_row[user]))
			current_row[user] = 0
		}
		// write it
		f.WriteString(fmt.Sprintf("%s, %s\n", current_hour, strings.Join(messages, ",")))
		// reset
		messages = nil
	}
	return nil
}

func (h Histogram) write_hourly_csv() (error) {
	const filename string = "./hourly.csv"
	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write csv file %s", filename))
	}

	defer f.Close()
	// write the header
	f.WriteString("Hour, MessageCount\n")

	for hour, _ := range h.Hours {
		f.WriteString(fmt.Sprintf("%02d:00, %d\n", hour, h.Hours[hour]))
	}

	return nil

}

func (h Histogram) write_weekday_csv() (error) {
	const filename string = "./weekday.csv"

	f, err := os.Create(filename)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write csv file %s", filename))
	}

	defer f.Close()

	// write the header
	f.WriteString("Weekday, MessageCount\n")

	for day, _ := range h.Weekdays {
		f.WriteString(fmt.Sprintf("%s, %d\n", day, h.Weekdays[day]))
	}

	return nil
}

func (h Histogram) get_chattiest_day() (time.Weekday, int) {

	var chattiest_day time.Weekday
	var chattiest_day_messages = 0

	for day, messages := range h.Weekdays {
		if messages > chattiest_day_messages {
			chattiest_day = day
			chattiest_day_messages = messages
		}
	}

	return chattiest_day, chattiest_day_messages
}

func (h Histogram) get_chattiest_hour() (int, int) {

	var chattiest_hour, chattiest_hour_messages = 0, 0

	for hour, messages := range h.Hours {
		if messages > chattiest_hour_messages {
			chattiest_hour = hour
			chattiest_hour_messages = messages
		}
	}

	return chattiest_hour, chattiest_hour_messages
}

func (h Histogram) report() {
	fmt.Printf("Total message sent: %d\n", h.TotalMessages)
	for user, messages := range h.Users {
		fmt.Printf("User %s sent %d messages in total.\n", strings.Trim(user, " "), messages)
	}

	day, daily_messages := h.get_chattiest_day()
	fmt.Printf("You send the most messages on %s, %d all told!\n", day, daily_messages)

	hour, hourly_messages := h.get_chattiest_hour()
	fmt.Printf("You send the most messages during the %dth hour of the day, %d all told!\n", hour, hourly_messages)


	h.write_weekday_csv()
	h.write_hourly_csv()
	h.write_alltime_csv()
}

func (m Message) display() {
	fmt.Println(m.Day, m.Month, m.Year, m.Hour, m.Minute, m.Second, m.User, m.Body)
	return
}

func Analyze() error {
	filename, err := parseArgs()

	if err != nil {
		return err
	}

	// Make sure the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Unable to open file: %s", filename))
	}

	fmt.Printf("Beginning analysis of %s ...\n", filename)

	f, err := os.Open(filename)

	if err != nil {
		return errors.New(fmt.Sprintf("Unable to access file: %s", filename))
	}

	// make sure the file gets closed when we exit
	defer f.Close()

	histo := new(Histogram)
	histo.init()

	scanner := bufio.NewScanner(f)
	line_count := 0
	for scanner.Scan() {
		line_count++
		message := new(Message)

		line := scanner.Text()

		// parse out the date and time
		fmt.Sscanf(line, "%d.%d.%d %d:%d:%d\n", &message.Day, &message.Month, &message.Year, &message.Hour, &message.Minute, &message.Second)

		separated := strings.SplitN(line, ",", 2)
		// looks like a safe assumption that correctly formatted records have a comma between the datestamp and the data
		if len(separated) < 2 {
			fmt.Printf("Warning: discarding malformed record on line: %d\n", line_count)
			fmt.Printf("\t%s\n", line)
			continue
		} else {
			// extract username and message body, again presence of colon seems like a safe assumption
			userbody := strings.SplitN(separated[1], ":", 2)
			if len(userbody) < 2 {
				fmt.Println("Warning: discarding malformed record.")
				fmt.Printf("\t%s\n", line)
				continue
			}
			message.User = userbody[0]
			message.Body = userbody[1]
		}
		histo.count(message)
	}

	histo.report()
	return nil
}
