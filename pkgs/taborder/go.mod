module github.com/KeithWiles/go-pktgen/pkgs/taborder

replace github.com/KeithWiles/go-pktgen/pkgs/ttylog => ../ttylog

go 1.18

require (
	github.com/KeithWiles/go-pktgen/pkgs/ttylog v0.0.0-00010101000000-000000000000
	github.com/gdamore/tcell/v2 v2.4.1-0.20210905002822-f057f0a857a1
)

require (
	github.com/gdamore/encoding v1.0.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2 // indirect
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d // indirect
	golang.org/x/text v0.3.6 // indirect
)
