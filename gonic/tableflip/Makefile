mockClient:
	curl -X GET http://localhost:8080/version
	while true; do curl http://localhost:8080/version;sleep 1; done

upgrade:
	# modify the version number in the main.go source file
	# and then run this target to update the version number in the
	# other files

	kill -s HUP $pid