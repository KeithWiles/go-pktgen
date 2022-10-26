module github.com/KeithWiles/go-pktgen/pkgs/devbind

replace github.com/KeithWiles/go-pktgen/pkgs/ttylog => ../ttylog

go 1.18

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/KeithWiles/go-pktgen/pkgs/ttylog v0.0.0-00010101000000-000000000000
	github.com/davecgh/go-spew v1.1.1
)
