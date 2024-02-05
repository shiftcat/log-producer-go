package util

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func ShutdownHook(f func()) chan bool {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		f()
		done <- true
	}()
	return done
}

func ToJson(v any) string {
	jsonb, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	return string(jsonb)
}
