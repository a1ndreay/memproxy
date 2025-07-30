package fake

import (
	"bytes"
	"testing"

	"github.com/a1ndreay/memproxy/pkg/cache"
)

func TestGet(t *testing.T) {
	tests := []struct {
		giveKey   string
		giveValue []byte
		want      []byte
		wantErr   error
	}{
		{
			giveKey:   "/index",
			giveValue: []byte(`<!DOCTYPE html>`),
			want:      []byte(`<!DOCTYPE html>`),
			wantErr:   nil,
		},
		{
			giveKey:   "missing",
			giveValue: nil,
			want:      nil,
			wantErr:   nil,
		},
	}

	var backend cache.Backend
	backend = New()

	for _, tt := range tests {
		t.Run(tt.giveKey, func(t *testing.T) {

			if tt.giveValue != nil {
				err := backend.Set(tt.giveKey, tt.giveValue)
				if err != nil {
					t.Fatalf("Set failed: %v", err)
				}
			}

			got, err := backend.Get(tt.giveKey)

			if err != tt.wantErr {
				t.Errorf("Get(%q) error = %v, wantErr %v", tt.giveKey, err, tt.wantErr)
			}

			if !bytes.Equal(got, tt.want) {
				t.Errorf("Get(%q) = %s, want %s", tt.giveKey, got, tt.want)
			}
		})
	}
}
