PORT?=40100

.PHONY: build
build:
	mvn clean install -DskipTests

.PHONY: run
run:
	java -jar ./target/evomaster.jar --sutControllerPort=${PORT} --outputFolder=/Users/juan.basso/Documents/University/thesis/EMB-golang/go/rest/artificial/ncs/generated_${lang}_${time} --maxTime=${time} --outputFormat=GO --seed=1