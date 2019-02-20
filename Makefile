install:
	go install
.PHONY: install

devstatic:
	go-bindata -debug -prefix 'static/' -o static.go static
