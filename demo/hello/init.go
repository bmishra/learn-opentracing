package hello

import (
	"context"
	"expvar"
	"log"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	oLog "github.com/opentracing/opentracing-go/log"
	"gopkg.in/tokopedia/logging.v1"
)

type ServerConfig struct {
	Name string
}

type Config struct {
	Server ServerConfig
}

type HelloWorldModule struct {
	cfg       *Config
	something string
	stats     *expvar.Int
}

func NewHelloWorldModule() *HelloWorldModule {

	var cfg Config

	ok := logging.ReadModuleConfig(&cfg, "config", "hello") || logging.ReadModuleConfig(&cfg, "files/etc/gosample", "hello")
	if !ok {
		// when the app is run with -e switch, this message will automatically be redirected to the log file specified
		log.Fatalln("failed to read config")
	}

	// this message only shows up if app is run with -debug option, so its great for debugging
	logging.Debug.Println("hello init called", cfg.Server.Name)

	return &HelloWorldModule{
		cfg:       &cfg,
		something: "John Doe",
		stats:     expvar.NewInt("rpsStats"),
	}

}

func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
	span, ctx := opentracing.StartSpanFromContext(r.Context(), r.URL.Path)
	defer span.Finish()

	hlm.stats.Add(1)
	hlm.someSlowFuncWeWantToTrace(ctx, w)
}

func (hlm *HelloWorldModule) someSlowFuncWeWantToTrace(ctx context.Context, w http.ResponseWriter) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "someSlowFuncWeWantToTrace")
	defer span.Finish()
	span.LogFields(
		oLog.String("slowFunc", "around 7 sec"),
	)

	//dummy time consuming functions
	callRedis(ctx)
	callDB(ctx)

	w.Write([]byte("Hello " + hlm.something))
}

func callRedis(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "callRedis")
	defer span.Finish()
	time.Sleep(time.Second * 3)
	return nil
}

func callDB(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "callDB")
	defer span.Finish()
	time.Sleep(time.Second * 4)
	return nil
}
