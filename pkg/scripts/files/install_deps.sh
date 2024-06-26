#!/usr/bin/env bash

BASE_DIR=$(dirname $(realpath -s $0))
source $BASE_DIR/common.sh

build_socat(){
    SOCAT_VERSION="1.7.3.4"
    local socat_tar="${BASE_DIR}/../components/socat-${SOCAT_VERSION}.tar.gz"

    if [ -f "$socat_tar" ]; then
        ensure_success $sh_c "cp ${socat_tar} ./"
    else
        ensure_success $sh_c "curl ${CURL_TRY} -LO http://www.dest-unreach.org/socat/download/socat-${SOCAT_VERSION}.tar.gz"
    fi
    # todo
    # ensure_success $sh_c "tar xzvf socat-${SOCAT_VERSION}.tar.gz"
    # ensure_success $sh_c "cd socat-${SOCAT_VERSION}"

    # ensure_success $sh_c "./configure --prefix=/usr && make -j4 && make install && strip socat"
}

build_contrack(){
    local contrack_tar="${BASE_DIR}/../components/conntrack-tools-1.4.1.tar.gz"
    if [ -f "$contrack_tar" ]; then
        ensure_success $sh_c "cp ${contrack_tar} ./"
    else
        ensure_success $sh_c "curl ${CURL_TRY} -LO https://github.com/fqrouter/conntrack-tools/archive/refs/tags/conntrack-tools-1.4.1.tar.gz"
    fi
    # todo
    # ensure_success $sh_c "tar zxvf conntrack-tools-1.4.1.tar.gz"
    # ensure_success $sh_c "cd conntrack-tools-1.4.1"

    # ensure_success $sh_c "./configure --prefix=/usr && make -j4 && make install"
}

install_deps() {
    case "$lsb_dist" in
        ubuntu|debian|raspbian)
            pre_reqs="apt-transport-https ca-certificates curl"
			if ! command -v gpg > /dev/null; then
				pre_reqs="$pre_reqs gnupg"
			fi
            ensure_success $sh_c 'apt-get update -qq >/dev/null'
            ensure_success $sh_c "DEBIAN_FRONTEND=noninteractive apt-get install -y -qq $pre_reqs >/dev/null"
            ensure_success $sh_c 'DEBIAN_FRONTEND=noninteractive apt-get install -y conntrack socat apache2-utils ntpdate net-tools make gcc openssh-server >/dev/null'
            ;;

        centos|fedora|rhel)
            if [ "$lsb_dist" = "fedora" ]; then
                pkg_manager="dnf"
            else
                pkg_manager="yum"
            fi

            ensure_success $sh_c "$pkg_manager install -y conntrack socat httpd-tools ntpdate net-tools make gcc openssh-server >/dev/null"
            ;;
        *)
            # build from source code
            # todo move to terminus-cli
            build_socat
            build_contrack

            #TODO: install bcrypt tools
            ;;
    esac
}


echo "---install_deps---"
install_deps
exit