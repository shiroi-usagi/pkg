package connconv

import "testing"

func TestParesClickhouseConnectionURL(t *testing.T) {
	type args struct {
		rawurl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "invalid url",
			args:    args{rawurl: "%"},
			wantErr: true,
		},
		{
			name: "without credentials",
			args: args{rawurl: "clickhouse://host:9999/db?debug=true"},
			want: "tcp://host:9999?database=db&debug=true",
		},
		{
			name: "without password",
			args: args{rawurl: "clickhouse://user@host:9999/db?debug=true"},
			want: "tcp://host:9999?database=db&debug=true&username=user",
		},
		{
			name: "with credentials",
			args: args{rawurl: "clickhouse://user:pass@host:9999/db?debug=true"},
			want: "tcp://host:9999?database=db&debug=true&password=pass&username=user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParesClickhouseConnectionURL(tt.args.rawurl)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertClickhouseConnectionURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("convertClickhouseConnectionURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParesMySQLConnectionURL(t *testing.T) {
	type args struct {
		rawurl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "invalid url",
			args:    args{rawurl: "%"},
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			args:    args{rawurl: "any://host:9999/db?debug=true"},
			wantErr: true,
		},
		{
			name: "without credentials",
			args: args{rawurl: "mysql://host:9999/db?debug=true"},
			want: "tcp(host:9999)/db?debug=true&parseTime=true",
		},
		{
			name: "without password",
			args: args{rawurl: "mysql://user@host:9999/db?debug=true"},
			want: "user@tcp(host:9999)/db?debug=true&parseTime=true",
		},
		{
			name: "with credentials",
			args: args{rawurl: "mysql://user:pass@host:9999/db?debug=true"},
			want: "user:pass@tcp(host:9999)/db?debug=true&parseTime=true",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParesMySQLConnectionURL(tt.args.rawurl)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParesMySQLConnectionURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParesMySQLConnectionURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
