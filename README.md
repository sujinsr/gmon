# gmon
Linux monitoring tool written in golang 

# Description
To build:

1. Setup golang

2. get termui source

    go get github.com/gizak/termui
    
3. To build the project

    make

To run

	./gmon <options>

    options
        m - Memory Info
        c - CPU Info(not yet implemented)
        d - Disk Info
        u - System uptime

To exit

	q

# License
GPL2, See License File
