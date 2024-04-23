PORT?=40100

.PHONY: build
build:
	@mvn clean install -DskipTests

.PHONY: run
run:
	@java -jar ./target/evomaster.jar --sutControllerPort=${PORT} --outputFolder=/Users/juan.basso/Documents/University/thesis/EMB-golang/go/rest/artificial/ncs/generated_${lang}_${time} --maxTime=${time} --outputFormat=GO --seed=1

OUTPUT_FORMAT?=GO
.PHONY: blackbox
blackbox:
	@java -jar ./target/evomaster.jar  --blackBox=true --bbSwaggerUrl=http://localhost:5001/swagger.json  --bbTargetUrl=http://localhost:5001 --maxTime=${time} --outputFormat=${OUTPUT_FORMAT} --ratePerMinute=60 --seed=1 --outputFolder=/Users/juan.basso/Documents/University/thesis/EMB-golang/go/rest/artificial/ncs/generated_bb_${lang}_${time}