package usecases

import (
	"time"

	"bitbucket.org/pborgeest/activity/entities"
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
