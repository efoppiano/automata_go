copy_env_file:
	if [ ! -f .env ]; then cp .env.example .env; fi

common: copy_env_file
	mkdir -p results

build:
	cd src && CGO_ENABLED=0 GOOS=linux go build -o automata.o
	mv src/automata.o .
	
scenario_1: build common
	GREEN_LIGHT_TIME=50 CROSSWALK_ROWS=6 \
	./automata.o

scenario_2: build common
	GREEN_LIGHT_TIME=35 CROSSWALK_ROWS=10 \
	./automata.o