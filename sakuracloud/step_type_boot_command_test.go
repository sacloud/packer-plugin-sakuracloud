package sakuracloud

import "testing"

func TestVNCAddress(t *testing.T) {
	tests := []struct {
		name string
		host string
		port string
		want string
	}{
		{
			name: "IPv4",
			host: "192.0.2.1",
			port: "5900",
			want: "192.0.2.1:5900",
		},
		{
			name: "IPv6",
			host: "2001:db8::1",
			port: "5900",
			want: "[2001:db8::1]:5900",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vncAddress(tt.host, tt.port); got != tt.want {
				t.Errorf("vncAddress(%q, %q) = %q, want %q", tt.host, tt.port, got, tt.want)
			}
		})
	}
}
