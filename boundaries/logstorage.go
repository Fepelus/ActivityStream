/*
Acts - add, display and delete activities to do.
Copyright (C) 2014  Patrick Borgeest
See LICENSE.txt for terms of usage.
*/

package boundaries

import (
	"crypto/sha1"
	"fmt"
	"os"
	"regexp"
	"time"

	"bufio"
	"github.com/Fepelus/ActivityStream/entities"
)

type Logfile struct {
	Filename string
}

type LogLine struct {
	Id       string
	Now      time.Time
	Command  string
	Activity entities.OneActivity
}

const (
	idxLength = 3
	Tformat   = "2006-01-02T15:04:05"
	Bformat   = "2006-01-02 15:04"
)

func (this LogLine) String() string {
	return fmt.Sprintf("[%s] %s\n", this.Id[0:idxLength], this.Activity)
}
func (this LogLine) LogString() string {
	now := this.Now.Format(Tformat)
	return fmt.Sprintf("[%s] %s: (%s) %s\n", now, this.Command, this.Id, this.Activity)
}

/* example input: "[2014-07-13T19:24:09] ADD: (414a4ec94c5b4c0f859b5f7cf721fceba05b4d84) 2014-05-05 05:07  Bam!" */
// ParseLogLine will take a single line of the logfile format
// and return the LogLine struct that represents it.
func ParseLogLine(input string) LogLine {
	regstring := "\\[(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2})\\] ([^:]+): \\(([0-9a-f]+)\\) (\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2} .*)"
	r, _ := regexp.Compile(regstring)
	match := r.FindStringSubmatch(input)
	/*
	   [1]: 2014-07-13T19:24:09
	   [2]: ADD
	   [3]: 414a4ec94c5b4c0f859b5f7cf721fceba05b4d84
	   [4]: 2014-05-05 05:07  Bam!
	*/
	if match == nil {
		fmt.Println("COULD NOT MATCH INPUT: ", input)
		return LogLine{}
	}
	nowstamp, _ := time.Parse(match[1], Tformat)
	activity, _ := entities.ParseOneActivity(match[4])
	activity.Id = match[3][0:idxLength]
	return LogLine{match[3], nowstamp, match[2], activity}
}

/* Probably should return the error rather than eat it */
func (this Logfile) AddNew(activity entities.OneActivity) string {

	thisLine := LogLine{
		sha(activity.String()),
		time.Now(),
		"ADD",
		activity,
	}

	_ = this.appendThisLine(thisLine)

	return thisLine.Id[0:idxLength]
}

func (this Logfile) appendThisLine(logline LogLine) error {
	f, err := os.OpenFile(this.Filename, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(this.Filename)
			if err == nil {
				f.Close()
				return this.appendThisLine(logline)
			}
		}
		return err
	}

	defer f.Close()
	_, err = f.WriteString(logline.LogString())
	return err
}

func sha(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (this Logfile) GetAll() entities.Activities {
	foundActivities := map[string]entities.OneActivity{}
	f, _ := os.Open(this.Filename)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		thisline := ParseLogLine(scanner.Text())
		if thisline.Command == "ADD" {
			foundActivities[thisline.Id] = thisline.Activity
		}
		if thisline.Command == "DELETE" {
			delete(foundActivities, thisline.Id)
		}
	}
	output := entities.Activities{}
	for _, v := range foundActivities {
		output = append(output, v)
	}
	return output
}

func (this Logfile) FindActivity(id string) entities.Activities {
	loglines := []LogLine{}
	output := entities.Activities{}
	f, _ := os.Open(this.Filename)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		thisline := ParseLogLine(scanner.Text())
		if thisline.Id[0:len(id)] == id {
			loglines = append(loglines, thisline)
		}
	}
	for _, el := range removeDeletedActivities(loglines) {
		output = append(output, el.Activity)
	}
	return output
}

func removeDeletedActivities(input []LogLine) []LogLine {
	output := []LogLine{}
	deletedList := []string{}
	for _, firstlogline := range input {
		if firstlogline.Command == "DELETE" {
			deletedList = append(deletedList, firstlogline.Id)
		}
	}
	for _, secondlogline := range input {
		if !stringInSlice(secondlogline.Id, deletedList) {
			output = append(output, secondlogline)
		}
	}
	return output
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (this Logfile) Delete(activity entities.OneActivity) error {
	thisLine := LogLine{
		sha(activity.String()),
		time.Now(),
		"DELETE",
		activity,
	}

	return this.appendThisLine(thisLine)
}
