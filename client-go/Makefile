FULL_PATH?=/Users/jcbasso/Documents/University/EvoMaster-golang/client-go
BIN?=${FULL_PATH}/bin
INSTRUMENTATION_BIN?=${BIN}/instrumentation
MAX_TIME?=3m

# EvoMaster
.PHONY: evomaster
evomaster:
	java -jar ./target/evomaster.jar --sutControllerPort=40101 --outputFolder=generated --maxTime=${MAX_TIME}

# Instrumentation
.PHONY: instrumentation
instrumentation:
	make clean
	@test -d ${BIN} || mkdir -p ${BIN}
	$(shell test -f $INSTRUMENTATION_BIN && rm $INSTRUMENTATION_BIN)
	go build -race -o=${INSTRUMENTATION_BIN} .

# Utils
.PHONY: clean
clean:
	go clean -cache