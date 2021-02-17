package main

import (
	"github.com/emirmuminoglu/emir"
)

func main() {
	e := emir.New(emir.Config{GracefulShutdown: true, Addr: ":8081"})

	e.GET("/", func(ctx emir.Context) error {
		ctx.Logger().Info("helloo")
		return nil
	})

	panic(e.ListenAndServe())
}
