.PHONY: clean all
SRC=boundaries/**/*go entities/**/*go usecases/**/*go 

all: acts server

acts: bin/cli/acts 
	cp $< .

server: bin/http/server
	cp $< $@

bin/cli/acts: bin/cli/acts.go ${SRC}
	cd $(<D); go build $(<F)

bin/http/server: bin/http/server.go ${SRC}
	cd $(<D); go build $(<F)

clean: 
	rm acts server bin/cli/acts bin/http/server
