#!/usr/bin/env bash

BASE_DIR=$(dirname $(realpath -s $0))
source $BASE_DIR/common.sh

local_ip=$1
os_platform=$2
ubuntuversion=$3

if [ -z "$os_platform" ]; then
  os_platform="$lsb_dist"
fi

precheck_os() {
    # try to resolv hostname
    ensure_success $sh_c "hostname -i >/dev/null"

    if [ -z "$local_ip" ]; then
        ip=$(ping -c 1 "$HOSTNAME" |awk -F '[()]' '/icmp_seq/{print $2}')
        printf "%s\t%s\n\n" "$ip" "$HOSTNAME"
        local_ip="$ip"
    fi
    

    # disable local dns
    case "$os_platform" in
        ubuntu|debian|raspbian)
            if system_service_active "systemd-resolved"; then
                ensure_success $sh_c "systemctl stop systemd-resolved.service >/dev/null"
                ensure_success $sh_c "systemctl disable systemd-resolved.service >/dev/null"
                if [ -e /usr/bin/systemd-resolve ]; then
                    ensure_success $sh_c "mv /usr/bin/systemd-resolve /usr/bin/systemd-resolve.bak >/dev/null"
                fi
                if [ -L /etc/resolv.conf ]; then
                    ensure_success $sh_c 'unlink /etc/resolv.conf && touch /etc/resolv.conf'
                fi
                config_resolv_conf
            else
                ensure_success $sh_c "cat /etc/resolv.conf > /etc/resolv.conf.bak"
            fi
            ;;
        centos|fedora|rhel)
            ;;
        *)
            ;;
    esac

    if ! hostname -i &>/dev/null; then
        ensure_success $sh_c "echo $local_ip  $HOSTNAME >> /etc/hosts"
    fi

    ensure_success $sh_c "hostname -i >/dev/null"

    # network and dns
    http_code=$(curl ${CURL_TRY} -sL -o /dev/null -w "%{http_code}" https://download.docker.com/linux/ubuntu)
    if [ "$http_code" != 200 ]; then
        config_resolv_conf
        if [ -f /etc/resolv.conf.bak ]; then
            ensure_success $sh_c "rm -rf /etc/resolv.conf.bak"
        fi
    fi

    # ubuntu 24 upgrade apparmor
    if [ ${ubuntuversion} = "2" ]; then
        aapv=$(apparmor_parser --version)
        if [[ ! ${aapv} =~ "4.0.1" ]]; then
            local aapv_tar="${BASE_DIR}/../components/apparmor_4.0.1-0ubuntu1_amd64.deb"
            if [ ! -f "$aapv_tar" ]; then
                ensure_success $sh_c "curl ${CURL_TRY} -k -sfLO --output-dir ${BASE_DIR}/../components/ https://launchpad.net/ubuntu/+source/apparmor/4.0.1-0ubuntu1/+build/28428840/+files/apparmor_4.0.1-0ubuntu1_amd64.deb"
            fi
            # todo test
            ensure_success $sh_c "echo 'apparmor setup'"
            # ensure_success $sh_c "dpkg -i ${BASE_DIR}/../components/apparmor_4.0.1-0ubuntu1_amd64.deb"
        fi
    fi

    # opy pre-installation dependency files 
    if [ -d /opt/deps ]; then
        ensure_success $sh_c "mv /opt/deps/* ${BASE_DIR}"
    fi
}

echo ">>> precheck_os.sh [$local_ip] [$os_platform] [$ubuntuversion]"
precheck_os
exit