/*
Acts - add, display and delete activities to do.
Copyright (C) 2014  Patrick Borgeest
See LICENSE.txt for terms of usage.
*/

package usecases

import (
	"fmt"

	"bytes"
	"github.com/Fepelus/ActivityStream/entities"
)

type CommandDeleter interface {
	FindActivity(id string) entities.Activities
	Delete(activity entities.OneActivity) error
}

/*
 * Basic flow :-
 * The user passes the ID.
 * The usecase fetches the single matching activity
 * The usecase gives the 'done' command to the deleter with this activity
 *
 * Alternative flows :-
 *  if the ID matches no activities then return a message to the user
 *  if the ID matches several activities then return them to the user and request a new ID
 */
func MarkActivityAsDone(id string, deleter CommandDeleter) error {
	activities := deleter.FindActivity(id)

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

	deleter.Delete(activities[0])
	return nil
}
