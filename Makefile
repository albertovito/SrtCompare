default:
	@GOOS=linux go build -o SrtCompare main.go

windows:
	@GOOS=windows go build -ldflags="-H windowsgui"

ico:
	@rsrc -ico ./res/srt.ico -o rsrc.syso