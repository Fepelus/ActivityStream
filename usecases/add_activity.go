package usecases

import "bitbucket.org/pborgeest/activity/entities"

type CommandAdder interface {
	AddNew(entities.OneActivity) string
}

// AddItem saves the given OneActivity in the datastorage passed as adder
// it returns a string that represents the ID of the new item in storage.
func AddItem(cmd entities.OneActivity, adder CommandAdder) string {
	return adder.AddNew(cmd)
}
