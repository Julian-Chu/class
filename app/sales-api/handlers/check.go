package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
)

func health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		return errors.New("application error")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}
	return json.NewEncoder(w).Encode(status)
}
