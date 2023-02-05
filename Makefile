test: 
	go run main.go > tests/tmp.txt
	diff tests/snapshot.txt tests/tmp.txt

update-snapshot:
	go run main.go > tests/snapshot.txt

clean:
	rm tests/tmp.txt
