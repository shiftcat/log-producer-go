module example.com

go 1.21

require (
	example.com/message v0.0.0
)

replace (
	example.com/message v0.0.0 => ./internal/message
)

require (
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
