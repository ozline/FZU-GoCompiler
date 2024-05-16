target:
	go build -o ./bin/GoLex
	make runcases

runcases:
	for i in {1..4}; do \
		./bin/GoLex < ./tests/case$$i.in > ./outs/case$$i.out; \
	done