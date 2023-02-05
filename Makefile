test: 
	go run main.go > tests/tokenize-tmp.txt
	diff tests/tokenize-output.txt tests/tokenize-tmp.txt

clean:
	rm tests/tokenize-tmp.txt

