package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

type Activities []OneActivity

type OneActivity struct {
	Timestamp  time.Time
	CommandTag string
	Body       string
}

func main() {
	activities := bufferedReaderToActivities(bufio.NewReader(os.Stdin))
	activities.sort()
	fmt.Print(activities.onlyThoseBefore(time.Now()))
}

func bufferedReaderToActivities(input *bufio.Reader) Activities {
	output := Activities{}
	for {
		line, err := input.ReadString('\n')
		if err != nil {
			break
		}
		activity, err := ParseOneActivity(line)
		if err == nil {
			output = append(output, activity)
		}
	}
	return output
}

/* Expected input:
 *    "YYYY-MM-DD HH:MM @rtask:every-n-hours:48 SRS a headline in Italian"
 * The "@rtask:" is optional
 * Output is a OneActivity struct with:
 *    - a go timestamp corresponding to YYYY-MM-DD MM:DD in Melbourne
 *    - everything between "@rtask:" and the next space character if the rtask
 *        is there ("" if it is not)
 *    - everything after the rtask if it is there otherwise everything after
 *        the timestamp
 * Throws an error if the timestamp cannot be parsed
 */
func ParseOneActivity(input string) (OneActivity, error) {
	loc, _ := time.LoadLocation("Australia/Melbourne")
	stamp, err := time.ParseInLocation("2006-01-02 15:04", input[0:16], loc)
	commandTag := ""
	body := ""
	if err != nil {
		return OneActivity{}, err
	}

	if len(input) > 23 && input[17:23] == "@rtask" {
		spaceIndex := strings.Index(input[24:len(input)], " ") + 24
		if spaceIndex > 23 {
			commandTag = input[24:spaceIndex]
			body = input[spaceIndex+1 : len(input)]
		} else {
			body = input[17:len(input)]
		}
	} else {
		body = input[17:len(input)]
	}

	return OneActivity{
		stamp,
		commandTag,
		body,
	}, nil
}

// Returns a new Activities struct containing only those
// in 'this' that have timestamps earlier than that given
/* Must preserve order */
func (this Activities) onlyThoseBefore(cutoff time.Time) Activities {
	for i := len(this); i > 0; i-- {
		if this[i-1].Timestamp.Before(cutoff) {
			return this[0:i]
		}
	}
	return Activities{}
}

func (this OneActivity) String() string {
	return fmt.Sprintf("%s  %s", this.Timestamp.Format("2006-01-02 15:04"), this.Body)
}

// Prints out each activity.
// First in the slice are printed last
func (this Activities) String() string {
	var buffer bytes.Buffer
	for i := len(this); i > 0; i-- {
		buffer.WriteString(this[i-1].String())
	}
	return buffer.String()
}

// Sorts in place. Earliest first
func (this Activities) sort() {
	sort.Sort(ByTime(this))
}

type ByTime Activities

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Timestamp.Before(a[j].Timestamp) }
