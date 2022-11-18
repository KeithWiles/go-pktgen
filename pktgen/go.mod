module github.com/KeithWiles/go-pktgen

replace github.com/KeithWiles/go-pktgen/pkgs/cfg => ../pkgs/cfg

replace github.com/KeithWiles/go-pktgen/pkgs/ttylog => ../pkgs/ttylog

replace github.com/KeithWiles/go-pktgen/pkgs/colorize => ../pkgs/colorize

replace github.com/KeithWiles/go-pktgen/pkgs/cpudata => ../pkgs/cpudata

replace github.com/KeithWiles/go-pktgen/pkgs/etimers => ../pkgs/etimers

replace github.com/KeithWiles/go-pktgen/pkgs/devbind => ../pkgs/devbind

replace github.com/KeithWiles/go-pktgen/pkgs/taborder => ../pkgs/taborder

replace github.com/KeithWiles/go-pktgen/pkgs/meter => ../pkgs/meter

go 1.19

require (
	github.com/KeithWiles/go-pktgen/pkgs/cfg v0.0.0-00010101000000-000000000000
	github.com/KeithWiles/go-pktgen/pkgs/colorize v0.0.0-20221026164806-7a528bb011d0
	github.com/KeithWiles/go-pktgen/pkgs/cpudata v0.0.0-20221026164806-7a528bb011d0
	github.com/KeithWiles/go-pktgen/pkgs/devbind v0.0.0-20221026164806-7a528bb011d0
	github.com/KeithWiles/go-pktgen/pkgs/etimers v0.0.0-20221026164806-7a528bb011d0
	github.com/KeithWiles/go-pktgen/pkgs/meter v0.0.0-00010101000000-000000000000
	github.com/KeithWiles/go-pktgen/pkgs/taborder v0.0.0-20221026164806-7a528bb011d0
	github.com/KeithWiles/go-pktgen/pkgs/ttylog v0.0.0-20221026164806-7a528bb011d0
	github.com/gdamore/tcell/v2 v2.5.3
	github.com/jessevdk/go-flags v1.5.0
	github.com/rivo/tview v0.0.0-20221117065207-09f052e6ca98
	github.com/shirou/gopsutil v3.21.11+incompatible
	golang.org/x/text v0.4.0
)

require (
	github.com/BurntSushi/toml v1.2.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/tidwall/jsonc v0.3.2 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.2.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
