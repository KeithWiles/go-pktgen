project('go-pktgen', 'C',
    version: run_command(find_program('cat', 'more'),
        files('VERSION'), check: false).stdout().strip(),

    license: 'BSD',
    default_options: [
        'buildtype=release',
        'default_library=static',
        'warning_level=3',
        'werror=true'
    ],
    meson_version: '>= 0.60.0'
)
   
go_pktgen_conf = configuration_data()

# set up some global vars for compiler, platform, configuration, etc.
cc = meson.get_compiler('c')

add_project_arguments('-D__PROJECT_VERSION="' + meson.project_version() + '"', language: 'c')
add_project_arguments('-D_GNU_SOURCE', language: 'c')

dpdk = dependency('libdpdk', required: true)
# message('prefix: ' + get_option('prefix') + ' libdir: ' + get_option('libdir'))

dpdk_libs_path = join_paths(get_option('prefix'), get_option('libdir'))
# message('DPDK lib path: ' + dpdk_libs_path)

dpdk_bond = cc.find_library('librte_net_bond', dirs: [dpdk_libs_path], required: false)

# enable extra warnings and disable any unwanted warnings
warning_flags = [
    '-Wno-pedantic',
    '-Wno-format-truncation',
]
foreach arg: warning_flags
    if cc.has_argument(arg)
        add_project_arguments(arg, language: 'c')
    endif
endforeach

subdir('libs')
subdir('pktgen')

