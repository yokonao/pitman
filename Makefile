SOURCE:=$(wildcard *.go)

test: test1 test2 test3 test4

test1: 
	go run $(SOURCE) -path=sample1 > tests/tmp1.txt
	diff tests/snapshot1.txt tests/tmp1.txt
test2:
	go run $(SOURCE) -path=sample2 > tests/tmp2.txt
	diff tests/snapshot2.txt tests/tmp2.txt
test3:
	go run $(SOURCE) -path=sample3 > tests/tmp3.txt
	diff tests/snapshot3.txt tests/tmp3.txt
test4:
	go run $(SOURCE) -path=partial_sample4 > tests/tmp4.txt
	diff tests/snapshot4.txt tests/tmp4.txt

update-snapshot:
	go run $(SOURCE) > tests/snapshot1.txt
	go run $(SOURCE) -path=sample2 > tests/snapshot2.txt
	go run $(SOURCE) -path=sample3 > tests/snapshot3.txt
	go run $(SOURCE) -path=partial_sample4 > tests/snapshot4.txt

clean:
	rm tests/tmp*.txt
