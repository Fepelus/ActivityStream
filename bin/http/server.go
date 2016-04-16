/*
Acts - add, display and delete activities to do.
Copyright (C) 2016  Patrick Borgeest
See LICENSE.txt for terms of usage.
*/

package main

import (
	"fmt"

	"bytes"
	"github.com/Fepelus/ActivityStream/boundaries"
	"github.com/Fepelus/ActivityStream/entities"
	"github.com/Fepelus/ActivityStream/usecases"
)

func main() {
	cmdToFunc := map[string]func([]string){
		"new":    newItem,
		"add":    newItem,
		"done":   doneItem,
		"del":    doneItem,
		"delete": doneItem,
		"delay":  delayItem,
		"grep":   grepItems,
		"get":    getActivity,
		"help":   help,
	}
	if len(os.Args) < 2 {
		help([]string{})
	} else if function, ok := cmdToFunc[os.Args[1]]; ok {
		function(os.Args[2:])
	} else {
		help([]string{})
	}
}

func getLogfile() boundaries.Logfile {
	if os.Getenv("ACTS_LOGFILE") != "" {
		return boundaries.Logfile{os.Getenv("ACTS_LOGFILE")}
	}
	return boundaries.Logfile{"logfile.txt"}
}

func newItem(args []string) {
	if doesNotPassTheCheck(args) {
		return
	}
	wholeInput := concatenate(args)

	activity, err := entities.ParseOneActivity(wholeInput)
	if err != nil {
		fmt.Println("You probably meant to say 'new now'")
		return
	}
	fmt.Println(usecases.AddItem(activity, getLogfile()))
}

func concatenate(args []string) string {
	var buffer bytes.Buffer
	for i := 0; i < len(args); i++ {
		if i > 0 {
			buffer.WriteString(" ")
		}
		buffer.WriteString(args[i])
	}
	return buffer.String()
}

func doesNotPassTheCheck(args []string) bool {
	if isNowCommand(args) {
		handleNowCommand(args)
		return true
	}
	if isTodayCommand(args) {
		handleTodayCommand(args)
		return true
	}
	if len(args) < 3 {
		fmt.Println("Arguments to newItem were only:", args)
		help(args)
		return true
	}
	return false
}

func isNowCommand(args []string) bool {
	return len(args) > 1 && args[0] == "now"
}
func handleNowCommand(args []string) {
	now := time.Now()
	newargs := todayWithTime(args, now.Format("15:04"))
	newItem(newargs)
}

func isTodayCommand(args []string) bool {
	_, err := time.Parse("15:04", args[0])
	return err == nil
}
func handleTodayCommand(args []string) {
	newargs := todayWithTime(args, args[0])
	newItem(newargs)
}

func todayWithTime(args []string, timestamp string) []string {
	now := time.Now()
	newargs := make([]string, len(args)+1)
	newargs[0] = now.Format("2006-01-02")
	newargs[1] = timestamp
	for i := 2; i < len(args)+1; i++ {
		newargs[i] = args[i-1]
	}
	return newargs
}

func getActivity(args []string) {
	items := usecases.GetActivity(getLogfile())
	for _, el := range items {
		fmt.Println(el)
	}
}

func doneItem(args []string) {
	if len(args) < 1 {
		fmt.Println("Arguments to doneItem were only:", args)
		help(args)
		return
	}
	err := usecases.MarkActivityAsDone(args[0], getLogfile())
	if err != nil {
		fmt.Print(err)
	}
}

func grepItems(args []string) {
	if len(args) < 1 {
		fmt.Println("Arguments to doneItem were only:", args)
		help(args)
		return
	}
	grepped := usecases.GrepItems("("+args[0], getLogfile())
	if grepped != nil {
		for _, el := range grepped {
			fmt.Println(el)
		}
	}
}

func delayItem(args []string) {
	//delay ID count unit ('hours' 'days')
	if len(args) < 1 {
		fmt.Println("Arguments to delayItem were only:", args)
		help(args)
		return
	}
	if len(args) == 1 {
		if err := usecases.DelayActivity(args[0], 1, "day", getLogfile()); err != nil {
			fmt.Print(err)
		}
		return
	}
	if len(args) == 2 {
		fmt.Println("Arguments to delayItem were only:", args)
		help(args)
		return
	}
	count, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Print(err)
	}
	if err = usecases.DelayActivity(args[0], count, args[2], getLogfile()); err != nil {
		fmt.Print(err)
	}
}

func parseActivity(datebit, timebit, body string) (entities.OneActivity, error) {
	loc, _ := time.LoadLocation("Australia/Melbourne")
	stamp, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", datebit, timebit), loc)
	if err != nil {
		return entities.OneActivity{}, err
	}
	return entities.OneActivity{"", stamp, "", body}, nil
}

func help(args []string) {
	fmt.Printf(`%s cmd [args]
    help
    get
    new [date] [time] [body]
    done [ID]
    grep [ID]
    delay [ID] [count] [unit] ('minutes' 'hours' 'days' 'weeks' 'months')
`, os.Args[0])
}
