
.PHONY: adapter advertising

adapter:
	BASEDIR=.. \
	FILE_FILTER=adapter \
	xAPI_FILTER="Provisioner" \
	xMETHOD_FILTER=Attach LOG_LEVEL=trace \
		go run ./srcgen/main.go full --debug

advertising:
	BASEDIR=.. \
	FILE_FILTER=advertising-api \
	xAPI_FILTER="Provisioner" \
	xMETHOD_FILTER=Attach LOG_LEVEL=trace \
	go run ./srcgen/main.go full --debug

mesh:
	BASEDIR=.. \
	FILE_FILTER=mesh-api \
	xAPI_FILTER="Provisioner" \
	xMETHOD_FILTER=Attach LOG_LEVEL=trace \
	go run ./srcgen/main.go full --debug

advertisement-monitor:
	BASEDIR=.. \
	FILE_FILTER=advertisement-monitor-api \
	xAPI_FILTER="Provisioner" \
	xMETHOD_FILTER=Attach LOG_LEVEL=trace \
	go run ./srcgen/main.go full --debug