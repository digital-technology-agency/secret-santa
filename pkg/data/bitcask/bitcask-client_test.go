package bitcask

import (
	"git.mills.io/prologic/bitcask"
	"reflect"
	"testing"
)

var testDbName = "test-db"

func TestConnect(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Data
		wantErr bool
	}{
		{
			name: "test table connect",
			args: args{
				testDbName,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Connect(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestData_Add(t *testing.T) {
	dataC, err := Connect(testDbName)
	if err != nil {
		t.Errorf("Connect() error = %v", err)
	}
	type fields struct {
		db *bitcask.Bitcask
	}
	type args struct {
		key   []byte
		value []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "add key to data",
			fields: fields{
				db: dataC.db,
			},
			args: args{
				key:   []byte("1"),
				value: []byte("значение"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Data{
				db: tt.fields.db,
			}
			if err := data.Add(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_Get(t *testing.T) {
	testKey := []byte("1")
	testValue := []byte("значение")
	dataC, err := Connect(testDbName)
	if err != nil {
		t.Errorf("Connect() error = %v", err)
	}
	err = dataC.Add(testKey, testValue)
	if err != nil {
		t.Errorf("Connect() error = %v", err)
	}
	type fields struct {
		db *bitcask.Bitcask
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "get data by key",
			fields: fields{
				db: dataC.db,
			},
			args: args{
				key: testKey,
			},
			want:    testValue,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Data{
				db: tt.fields.db,
			}
			got, err := data.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_GetAll(t *testing.T) {
	testKey := []byte("1")
	testValue := []byte("значение")
	dataC, err := Connect(testDbName)
	if err != nil {
		t.Errorf("Connect() error = %v", err)
	}
	err = dataC.Add(testKey, testValue)
	if err != nil {
		t.Errorf("Connect() error = %v", err)
	}
	type fields struct {
		db *bitcask.Bitcask
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string][]byte
		wantErr bool
	}{
		{
			name: "get all data",
			fields: fields{
				db: dataC.db,
			},
			want: map[string][]byte{
				"1": testValue,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Data{
				db: tt.fields.db,
			}
			got, err := data.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestData_Remove(t *testing.T) {
	testKey := []byte("1")
	testValue := []byte("значение")
	dataC, err := Connect(testDbName)
	if err != nil {
		t.Errorf("Connect() error = %v", err)
	}
	err = dataC.Add(testKey, testValue)
	if err != nil {
		t.Errorf("Connect() error = %v", err)
	}
	type fields struct {
		db *bitcask.Bitcask
	}
	type args struct {
		key []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "remove key",
			fields: fields{
				db: dataC.db,
			},
			args: args{
				key: testKey,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Data{
				db: tt.fields.db,
			}
			if err := data.Remove(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
