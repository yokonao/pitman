SOURCE:=$(wildcard *.go)


test1: 
	go run $(SOURCE) -path=sample1 > tests/tmp1.txt
	diff tests/snapshot1.txt tests/tmp1.txt
test2:
	go run $(SOURCE) -path=sample2 > tests/tmp2.txt
	diff tests/snapshot2.txt tests/tmp2.txt


update-snapshot:
	go run $(SOURCE) > tests/snapshot1.txt
	go run $(SOURCE) -path=sample2 > tests/snapshot2.txt

clean:
	rm tests/tmp.txt
