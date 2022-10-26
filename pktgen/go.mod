module github.com/KeithWiles/go-pktgen

replace github.com/KeithWiles/go-pktgen/pkgs/ttylog => ../pkgs/ttylog

replace github.com/KeithWiles/go-pktgen/pkgs/colorize => ../pkgs/colorize

replace github.com/KeithWiles/go-pktgen/pkgs/cpudata => ../pkgs/cpudata

replace github.com/KeithWiles/go-pktgen/pkgs/etimers => ../pkgs/etimers

replace github.com/KeithWiles/go-pktgen/pkgs/devbind => ../pkgs/devbind

replace github.com/KeithWiles/go-pktgen/pkgs/taborder => ../pkgs/taborder

go 1.19

require (
	github.com/KeithWiles/go-pktgen/pkgs/colorize v0.0.0-00010101000000-000000000000
	github.com/KeithWiles/go-pktgen/pkgs/cpudata v0.0.0-00010101000000-000000000000
	github.com/KeithWiles/go-pktgen/pkgs/devbind v0.0.0-00010101000000-000000000000
	github.com/KeithWiles/go-pktgen/pkgs/etimers v0.0.0-00010101000000-000000000000
	github.com/KeithWiles/go-pktgen/pkgs/taborder v0.0.0-20221022153719-9c0f739786b4
	github.com/KeithWiles/go-pktgen/pkgs/ttylog v0.0.0-00010101000000-000000000000
	github.com/gdamore/tcell/v2 v2.4.1-0.20210905002822-f057f0a857a1
	github.com/jessevdk/go-flags v1.5.0
	github.com/rivo/tview v0.0.0-20220916081518-2e69b7385a37
	github.com/shirou/gopsutil v3.21.11+incompatible
	golang.org/x/text v0.3.7
)

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
