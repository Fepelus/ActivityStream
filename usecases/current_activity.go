/*
Acts - add, display and delete activities to do.
Copyright (C) 2014  Patrick Borgeest
See LICENSE.txt for terms of usage.
*/

package usecases

import (
	"time"

	"github.com/Fepelus/ActivityStream/entities"
)

type CommandGetter interface {
	GetAll() entities.Activities
}

func GetActivity(getter CommandGetter) []string {
	output := []string{}
	activities := getter.GetAll()

	// order by user-entered datestamp
	activities.Sort()

	// only return those before now
	now := time.Now()
	for _, oneActivity := range activities {
		if now.Before(oneActivity.Timestamp) {
			break
		}
		output = append(output, oneActivity.IndexedString())
	}
	return output
}
