target:
	go build -o ./bin/GoParser
	make runcases

runcases:
	for i in {1..7}; do \
		./bin/GoParser < ./tests/case$$i.in > ./outs/case$$i.out; \
	done

test:
	go build -o ./bin/GoParser
	./bin/GoParser < ./tests/case5.in > ./outs/case5.out; \