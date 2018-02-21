package main

import (
	"context"
	"fmt"

	"time"

	"gopkg.in/tokopedia/logging.v1"
	"gopkg.in/tokopedia/logging.v1/tracer"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/tokopedia/learn-opentracing/config"
)

func main() {
	logging.LogInit()
	tracer.Init(&config.Config)

	ctx := context.Background()
	span, ctx := opentracing.StartSpanFromContext(ctx, "Demo App")
	defer span.Finish()

	a, b := getNum(ctx)
	s := sum(ctx, a, b)
	fmt.Println("Sum: ", s)
	fmt.Println(span)

	time.Sleep(time.Second * 20)
}

func getNum(ctx context.Context) (int, int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "getNum")
	defer span.Finish()
	return 2, 4
}

func sum(ctx context.Context, a, b int) int {
	span, ctx := opentracing.StartSpanFromContext(ctx, "sum")
	defer span.Finish()

	if a == 0 || b == 0 {
		ext.Error.Set(span, true)
	}
	return a + b
}
