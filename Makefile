target:
	go build -o ./bin/GoParser
	make runcases

runcases:
	for i in {1..7}; do \
		./bin/GoParser ok < ./tests/case$$i.in > ./outs/case$$i.out; \
	done

test:
	go build -o ./bin/GoParser
	./bin/GoParser ok < ./tests/case5.in > ./outs/case5.out; \