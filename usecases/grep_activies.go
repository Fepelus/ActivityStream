/*
Acts - add, display and delete activities to do.
Copyright (C) 2015  Patrick Borgeest
See LICENSE.txt for terms of usage.
*/

package usecases

type CommandGrepper interface {
	Grep(id string) []string
}

func GrepItems(id string, grepper CommandGrepper) []string {
	return grepper.Grep(id)
}
