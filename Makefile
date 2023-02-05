SOURCE:=$(wildcard *.go)


test: 
	go run $(SOURCE) > tests/tmp.txt
	diff tests/snapshot.txt tests/tmp.txt

update-snapshot:
	go run $(SOURCE) > tests/snapshot.txt

clean:
	rm tests/tmp.txt
