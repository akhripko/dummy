package service

import (
	"errors"
	"reflect"
	"testing"

	"github.com/akhripko/dummy/src/models"
	mock "github.com/akhripko/dummy/src/service/mock"
	"github.com/golang/mock/gomock"
)

func TestService_Hello(t *testing.T) {
	type fields struct {
		storage Storage
		cache   Cache
	}
	type args struct {
		name string
	}

	c := gomock.NewController(t)
	storage := mock.NewMockStorage(c)
	storage.EXPECT().Hello("key2").DoAndReturn(func(name string) (*models.HelloMessage, error) {
		return &models.HelloMessage{Message: "Hello, key2"}, nil
	}).Times(1)
	storage.EXPECT().Hello("key3").DoAndReturn(func(name string) (*models.HelloMessage, error) {
		return &models.HelloMessage{Message: "Hello, key3"}, nil
	}).Times(1)
	storage.EXPECT().Hello("key4").DoAndReturn(func(name string) (*models.HelloMessage, error) {
		return nil, errors.New("some error")
	}).Times(1)

	cache := mock.NewMockCache(c)
	cache.EXPECT().Read("key").DoAndReturn(func(key string) (string, error) { return "Hello, key", nil }).Times(1)
	cache.EXPECT().Read("key2").DoAndReturn(func(key string) (string, error) { return "", nil }).Times(1)
	cache.EXPECT().WriteTTL("key2", "Hello, key2", 300).DoAndReturn(func(key, value string, ttl int) error { return nil }).Times(1)
	cache.EXPECT().Read("key3").DoAndReturn(func(key string) (string, error) { return "", errors.New("some error") }).Times(1)
	cache.EXPECT().WriteTTL("key3", "Hello, key3", 300).DoAndReturn(func(key, value string, ttl int) error { return nil }).Times(1)
	cache.EXPECT().Read("key4").DoAndReturn(func(key string) (string, error) { return "", nil }).Times(1)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.HelloMessage
		wantErr bool
	}{
		{
			name:    "read_cached",
			fields:  fields{storage: storage, cache: cache},
			args:    args{name: "key"},
			want:    &models.HelloMessage{Message: "Hello, key"},
			wantErr: false,
		},
		{
			name:    "cache_empty | read_storage | save_cache",
			fields:  fields{storage: storage, cache: cache},
			args:    args{name: "key2"},
			want:    &models.HelloMessage{Message: "Hello, key2"},
			wantErr: false,
		},
		{
			name:    "cache_err | read_storage | save_cache",
			fields:  fields{storage: storage, cache: cache},
			args:    args{name: "key3"},
			want:    &models.HelloMessage{Message: "Hello, key3"},
			wantErr: false,
		},
		{
			name:    "cache_empty | storage_err",
			fields:  fields{storage: storage, cache: cache},
			args:    args{name: "key4"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
				cache:   tt.fields.cache,
			}
			got, err := s.Hello(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hello() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hello() got = %v, want %v", got, tt.want)
			}
		})
	}
}
