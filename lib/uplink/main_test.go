// Copyright (C) 2018 Storj Labs, Inc.
// See LICENSE for copying information.

package uplink

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"storj.io/storj/pkg/identity"
)

func TestNewUplink(t *testing.T) {
	type args struct {
		ident         *identity.FullIdentity
		satelliteAddr string
		cfg           Config
	}

	ctx := context.Background()
	testIdentity, err := identity.NewFullIdentity(ctx, 12, 4)
	assert.NoError(t, err)

	tests := []struct {
		name string
		args args
		want *Uplink
	}{
		{
			name: "Test LibUplink Default",
			args: args{
				ident:         testIdentity,
				satelliteAddr: "satellite.example.com",
				cfg:           Config{},
			},
			want: &Uplink{
				ID:            testIdentity,
				SatelliteAddr: "satellite.example.com",
				Config:        Config{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUplink(tt.args.ident, tt.args.satelliteAddr, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUplink() = %v, want %v", got, tt.want)
			}
		})
	}
}