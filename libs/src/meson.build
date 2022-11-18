
sources = files(
)

libpkt = library('pkt', sources, dependencies: [common, dpdk])
pkt = declare_dependency(link_with: libpkt, include_directories: include_directories('.'))

