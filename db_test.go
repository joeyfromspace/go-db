package db

import (
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

func TestInitializeValidation(t *testing.T) {
	type args struct {
		o *ConnectOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *mongo.Database
		wantErr bool
	}{
		0: {
			name: "no connection url",
			args: args{
				o: &ConnectOptions{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Initialize(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitialize(t *testing.T) {
	type args struct {
		o *ConnectOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *mongo.Database
		wantErr bool
	}{
		0: {
			name: "with connection url",
			args: args{
				o: &ConnectOptions{
					DatabaseURL: "mongodb://127.0.0.1:27017",
					Database:    "test",
				},
			},
			want:    &mongo.Database{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Initialize(tt.args.o)
			if (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("Initialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDatabase(t *testing.T) {
	tests := []struct {
		before  func()
		name    string
		want    *mongo.Database
		wantErr bool
	}{
		0: {
			before:  Close,
			name:    "error on unitialized database connection",
			want:    nil,
			wantErr: true,
		},
		1: {
			before: func() {
				Initialize(&ConnectOptions{
					DatabaseURL: "mongodb://localhost:27017",
					Database:    "test",
				})
			},
			name:    "get database instance",
			want:    &mongo.Database{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before()
			}
			got, err := GetDatabase()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("GetDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_selectDatabase(t *testing.T) {
	type args struct {
		c *mongo.Client
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    *mongo.Database
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := selectDatabase(tt.args.c, tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("selectDatabase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("selectDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}
