package iceflake

import (
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/istyle-inc/iceflake/foundation"
	"github.com/istyle-inc/iceflake/tests/mocks"
)

func TestIceFlakeGenerator_Generate(t *testing.T) {
	// setup mock
	c := gomock.NewController(t)
	defer c.Finish()
	defer func() {
		foundation.InternalTimer = foundation.NewLocalTimer()
	}()
	mock := mocks.NewMockTimer(c)
	mock.EXPECT().Now().Times(3).Return(
		time.Date(2018, 4, 1, 0, 0, 0, 0, time.UTC))
	foundation.InternalTimer = mock

	type fields struct {
		w        uint64
		baseTime time.Time
		lastTS   uint64
		seq      uint64
	}
	tests := []struct {
		name    string
		fields  fields
		want    uint64
		wantErr bool
	}{
		{
			name: "normal case",
			fields: fields{
				w:        1,
				baseTime: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
				lastTS:   0,
				seq:      initialSequentialNumber,
			},
			want:    32614907904004097,
			wantErr: false,
		},
		{
			name: "same time called",
			fields: fields{
				w:        1,
				baseTime: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
				lastTS:   7776000000,
				seq:      initialSequentialNumber,
			},
			want:    32614907904004098,
			wantErr: false,
		},
		{
			name: "called from oldtime",
			fields: fields{
				w:        1,
				baseTime: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
				lastTS:   7776000001,
				seq:      initialSequentialNumber,
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := IceFlakeGenerator{
				w:        tt.fields.w,
				baseTime: tt.fields.baseTime,
				lastTS:   tt.fields.lastTS,
				seq:      tt.fields.seq,
			}
			got, err := g.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("IceFlakeGenerator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IceFlakeGenerator.Generate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewIDGenerator(t *testing.T) {
	type args struct {
		workerID uint64
		baseTime time.Time
	}
	tests := []struct {
		name string
		args args
		want IDGenerator
	}{
		{
			name: "succeeded to make new generator",
			args: args{
				workerID: 1,
				baseTime: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: &IceFlakeGenerator{
				w:        1,
				baseTime: time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC),
				lastTS:   0,
				seq:      initialSequentialNumber,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewIDGenerator(tt.args.workerID, tt.args.baseTime); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewIDGenerator() = %v, want %v", got, tt.want)
			}
		})
	}
}
