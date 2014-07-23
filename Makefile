all: activity acts

#hud: hud.go
#	go build hud.go

activity: activity.go
	go build $<

acts: acts.go */*go
	go build $<


