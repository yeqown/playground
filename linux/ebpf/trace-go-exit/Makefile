.PHNOY: binary
binary: clear
	go build -gcflags="-N -l" -o binary ./main.go

trypanic:
	sudo bpftrace -c "./binary -mode panic" trace.bt

tryexit: 
	sudo bpftrace -c "./binary -mode exit -code 2" trace.bt

listuprobes:
	bpftrace -l "uprobe:./binary:*"

clear:
	rm ./binary