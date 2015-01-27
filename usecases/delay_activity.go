/*
Acts - add, display and delete activities to do.
Copyright (C) 2014  Patrick Borgeest
See LICENSE.txt for terms of usage.
*/

package usecases

import (
	"fmt"
	"time"

	"bytes"
	"github.com/Fepelus/ActivityStream/entities"
)

type CommandDelayer interface {
	CommandDeleter
	CommandAdder
}

//
// Basic flow :-
// The user passes the ID.
// The usecase fetches the single matching activity
// It gives the 'done' command to the delayer with this activity
// It creates a new command with the same details as the old command
// It alters the timestamp of the new command
// It sends the new command to the delayer to store as a new command
//
// Alternative flows :-
//  if the ID matches no activities then return a message to the user
//  if the ID matches several activities then return them to the user and request a new ID
//  if the unit is not among the legal strings then return a message to the user
//  In any of the alternative flows, no entries are written to the delayer.
//
func DelayActivity(id string, count int, unit string, delayer CommandDelayer) error {
	activities := delayer.FindActivity(id)

	if len(activities) == 0 {
		return fmt.Errorf("No activities found with index %s\n", id)
	}
	if len(activities) > 1 {
		var buffer bytes.Buffer
		buffer.WriteString("Ambiguous ID matches:\n")
		for i := 0; i < len(activities); i++ {
			buffer.WriteString(activities[i].IndexedString())
			buffer.WriteString("\n")
		}
		buffer.WriteString("\nNothing has been deleted. You may try again.\n")
		return fmt.Errorf(buffer.String())
	}
	thisActivity := activities[0]

	newtimestamp, err := delayTimestamp(thisActivity.Timestamp, count, unit)
	if err != nil {
		return err
	}
	delayer.Delete(activities[0])
	delayer.AddNew(entities.OneActivity{
		"",
		newtimestamp,
		thisActivity.CommandTag,
		thisActivity.Body,
	})
	return nil
}

func DelayActivityOneDay(id string, delayer CommandDelayer) error {
	return DelayActivity(id, 1, "day", delayer)
}

func delayTimestamp(input time.Time, count int, unit string) (time.Time, error) {
	if unit == "month" || unit == "months" {
		return input.AddDate(0, count, 0), nil
	}
	if unit == "week" || unit == "weeks" {
		return input.AddDate(0, 0, count*7), nil
	}
	if unit == "day" || unit == "days" {
		return input.AddDate(0, 0, count), nil
	}
	if unit == "hour" || unit == "hours" {
		return input.Add(time.Duration(count) * time.Hour), nil
	}
	if unit == "minute" || unit == "minutes" {
		return input.Add(time.Duration(count) * time.Minute), nil
	}

	return input, fmt.Errorf("Unit '%s' not found. Legal units are 'month','week','day','hour','minute'", unit)
}
