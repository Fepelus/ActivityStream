/*
Acts - add, display and delete activities to do.
Copyright (C) 2014  Patrick Borgeest
See LICENSE.txt for terms of usage.
*/

package entities

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"
)

type OneActivity struct {
	Id         string
	Timestamp  time.Time
	CommandTag string
	Body       string
}

func (this OneActivity) IndexedString() string {
	star := ""
	if this.HasRepeatCommand() {
		star = "*"
	}
	return fmt.Sprintf("[\033[1m%s\033[0m]%s %s %s", this.Id, star, this.TimeString(), this.Body)
}

func (this OneActivity) String() string {
	return fmt.Sprintf("%s %s", this.TimeString(), this.Body)
}

// Antipattern: Having code for something you aren't yet implementing...
func (this OneActivity) HasRepeatCommand() bool {
	return this.CommandTag != ""
}

func (this OneActivity) TimeString() string {
	return this.Timestamp.Format("2006-01-02 15:04")
}

func (this OneActivity) FullString() string {
	return fmt.Sprintf("%s %s %s", this.TimeString(), this.CommandTag, this.Body)
}

// Expected input:
//    "YYYY-MM-DD HH:MM @rtask:every-n-hours:48 SRS a headline in Italian"
// The "@rtask:" is optional
// Output is a OneActivity struct with:
//    - a go timestamp corresponding to YYYY-MM-DD MM:DD in Melbourne
//    - everything between "@rtask:" and the next space character if the rtask
//        is there ("" if it is not)
//    - everything after the rtask if it is there otherwise everything after
//        the timestamp
// Throws an error if the timestamp cannot be parsed
//
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
		"",
		stamp,
		commandTag,
		body,
	}, nil
}

type Activities []OneActivity

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
func (this Activities) Sort() {
	sort.Sort(ByTime(this))
}

type ByTime Activities

func (a ByTime) Len() int      { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool {
	if a[i].Timestamp.Before(a[j].Timestamp) {
		return true
	}
	if a[j].Timestamp.Before(a[i].Timestamp) {
		return false
	}
	return a[i].Body < a[j].Body
}
