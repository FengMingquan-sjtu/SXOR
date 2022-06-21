## Sequential XOR code
Project 1 of Big Data Processing Technologies 2022

Mingquan Feng 0210339100

## How to run
To check the performance of SXOR code:
```
go test -bench .
```

To compile and run MINIO with SXOR (the size MINIO source code is too large, please use git to download):
```
git clone https://github.com/FengMingquan-sjtu/minio
cd minio
go build -o minio-proj
./minio-proj server [your servers]
```