build:
	cd src && CGO_ENABLED=0 GOOS=linux go build -o automata.o
	mv src/automata.o .

run: build
	./automata.o