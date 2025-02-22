default:
	@GOOS=windows go build -ldflags="-H windowsgui" 
	@flatpak run org.winehq.Wine SrtComparator.exe