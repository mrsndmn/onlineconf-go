module github.com/mrsndmn/onlineconf-go

go 1.8

require (
	github.com/alldroll/cdb v1.0.3
	github.com/fsnotify/fsnotify v1.4.7
	github.com/grandecola/mmap v0.6.0
	github.com/sergi/go-diff v1.1.0
	github.com/stretchr/testify v1.6.1
	goa.design/goa v2.2.5+incompatible
	goa.design/goa/v3 v3.2.6
	golang.org/x/tools v0.0.0-20201229013931-929a8494cf60
)

replace github.com/alldroll/cdb => github.com/mrsndmn/cdb v1.0.3-0.20200206115956-e959e50d19c9
