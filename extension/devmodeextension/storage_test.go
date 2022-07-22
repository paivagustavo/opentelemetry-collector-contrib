package devmode

import (
	"context"
	"fmt"
	"testing"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()

	client, err := newClient(ctx, "sqlite3", "spans")

	if err != nil {
		fmt.Println(err.Error())
	}

	client.Set(ctx, "key", []byte("value"))

	v, _ := client.Get(ctx, "key")
	fmt.Println(string(v))
}
