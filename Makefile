all: acts

acts: acts.go */*go
	go build $<


