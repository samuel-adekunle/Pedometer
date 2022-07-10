#!/bin/bash

cd "$(dirname "$0")"
cd steps

for func in */; do
	cd ${func::-1}
	go build && zip ${func::-1}.zip ${func::-1}
	aws lambda update-function-code --function-name steps_${func::-1} --zip-file fileb://${func::-1}.zip > ${func::-1}.log
	cd ..
done

cd ../challenges

for func in */; do
	cd ${func::-1}
	go build && zip ${func::-1}.zip ${func::-1}
	aws lambda update-function-code --function-name challenges_${func::-1} --zip-file fileb://${func::-1}.zip > ${func::-1}.log
	cd ..
done