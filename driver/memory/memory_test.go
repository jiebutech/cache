package memory

import (
	"context"
	"github.com/jiebutech/cache"
	"github.com/jiebutech/config"
	"testing"
	"time"
)

var cli cache.Cache
var ctx = context.Background()

func TestMain(m *testing.M) {
	cli = NewMemoryCache(config.CacheConfig{})
	m.Run()

}

type A struct {
	B string `json:"b"`
	C int    `json:"c"`
}

func TestSet(t *testing.T) {
	key, value := "hello", "world"
	cli.Set(ctx, key, value, 3*time.Second)
	t.Log(cli.Get(ctx, key))
	time.Sleep(5 * time.Second)
	t.Log(cli.Get(ctx, key))
	cli.Set(ctx, "json", &A{
		B: "1231232131",
		C: 73982483420,
	})
	a := &A{}
	err := cli.GetAndUnmarshal(ctx, "json", a)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", a)
	cli.Clear(ctx)
	t.Log(cli.Get(ctx, key))
	err = cli.GetAndUnmarshal(ctx, "json", a)
	if err != nil {
		t.Fatal(err)
	}
}
