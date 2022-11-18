
sources = files(
    'pktgen-arp.c',
    'pktgen.c',
    'pktgen-capture.c',
    'pktgen-cfg.c',
    'pktgen-cmds.c',
    'pktgen-ether.c',
    'pktgen-gre.c',
    'pktgen-gtpu.c',
    'pktgen-ipv4.c',
    'pktgen-ipv6.c',
    'pktgen-latency.c',
    'pktgen-main.c',
    'pktgen-pcap.c',
    'pktgen-port-cfg.c',
    'pktgen-random.c',
    'pktgen-range.c',
    'pktgen-rate.c',
    'pktgen-seq.c',
    'pktgen-stats.c',
    'pktgen-tcp.c',
    'pktgen-udp.c',
    'pktgen-vlan.c',
)

libpktgen = library('pktgen', sources, dependencies: [common, dpdk])
pktgen = declare_dependency(link_with: libpktgen, include_directories: include_directories('.'))

