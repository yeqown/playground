# README

### 查询可追踪的uprobe函数列表，通配符过滤

sudo bpftrace -l 'uprobe:./binary:*'

```sh
# compile binary executable from go source file
go build -o binary ./main

# running bpftrace with trace.bt
# bpftrace -c "./binary -mode exit -code 3" trace.bt
bpftrace -c "./binary -mode panic" trace.bt
```