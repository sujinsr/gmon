all: gmon

SRC =	gmon.go\
		cpustat.go\
		diskstat.go\
		uptime.go\
		memstat.go\

gmon: $(SRC)
	go fmt $^
	go build -o gmon $^

clean:
	rm gmon
