if [ $1 = "run" ]
then
	if [ $2 = "appened" ]
		then
			docker stop appened
			docker rm appened
			docker run -d -v $PWD/data:/data --env-file .env --name appened -p 8082:8081 appened:latest
		elif [ $2 = "appened-twilio" ]
		then
			docker stop appened-twilio
			docker rm appened-twilio
			docker run -d --name appened-twilio -p 8083:8080 appended-twilio:latest
		else
			echo "Unknown service $2"
	fi
elif [ $1 = "build" ]
then
	if [ $2 = "appened" ]
	then
		docker build ./ -t appened:latest
	elif [ $2 = "appened-twilio" ]
	then
		docker build -f ./clients/twilio/Dockerfile -t appended-twilio:latest .
	else
		echo "Unknown service $2"
	fi
else
	echo "Unknown command $1"
fi
