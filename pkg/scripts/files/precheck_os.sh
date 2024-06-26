#!/usr/bin/env bash

BASE_DIR=$(dirname $(realpath -s $0))
source $BASE_DIR/common.sh

is_debian() {
    lsb_release=$(lsb_release -d 2>&1 | awk -F'\t' '{print $2}')
    if [ -z "$lsb_release" ]; then
        echo 0
        return
    fi
    if [[ ${lsb_release} == *Debian*} ]]; then
        case "$lsb_release" in
            *12.* | *11.*)
                echo 1
                ;;
            *)
                echo 0
                ;;
        esac
    else
        echo 0
    fi
}

is_ubuntu() {
    lsb_release=$(lsb_release -d 2>&1 | awk -F'\t' '{print $2}')
    if [ -z "$lsb_release" ]; then
        echo 0
        return
    fi
    if [[ ${lsb_release} == *Ubuntu* ]];then 
        case "$lsb_release" in
            *24.*)
                echo 2
                ;;
            *22.* | *20.*)
                echo 1
                ;;
            *)
                echo 0
                ;;
        esac
    else
        echo 0
    fi
}

precheck_os() {
    local ip os_type os_arch

    # check os type and arch and os vesion
    os_type=$(uname -s)
    os_arch=$(uname -m)
    os_verion=$(lsb_release -d 2>&1 | awk -F'\t' '{print $2}')

    if [ x"${os_type}" != x"Linux" ]; then
        log_fatal "unsupported os type '${os_type}', only supported 'Linux' operating system"
    fi

    if [[ x"${os_arch}" != x"x86_64" && x"${os_arch}" != x"amd64" ]]; then
        log_fatal "unsupported os arch '${os_arch}', only supported 'x86_64' architecture"
    fi

    if [[ $(is_ubuntu) -eq 0 && $(is_debian) -eq 0 ]]; then
        log_fatal "unsupported os version '${os_verion}', only supported Ubuntu 20.x, 22.x, 24.x and Debian 11, 12"
    fi

    # try to resolv hostname
    ensure_success $sh_c "hostname -i >/dev/null"

    ip=$(ping -c 1 "$HOSTNAME" |awk -F '[()]' '/icmp_seq/{print $2}')
    printf "%s\t%s\n\n" "$ip" "$HOSTNAME"

    local_ip="$ip"

    # disable local dns
    case "$lsb_dist" in
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
    ubuntuversion=$(is_ubuntu)
    if [ ${ubuntuversion} -eq 2 ]; then
        aapv=$(apparmor_parser --version)
        if [[ ! ${aapv} =~ "4.0.1" ]]; then
            local aapv_tar="${BASE_DIR}/../components/apparmor_4.0.1-0ubuntu1_amd64.deb"
            if [ ! -f "$aapv_tar" ]; then
                ensure_success $sh_c "curl ${CURL_TRY} -k -sfLO https://launchpad.net/ubuntu/+source/apparmor/4.0.1-0ubuntu1/+build/28428840/+files/apparmor_4.0.1-0ubuntu1_amd64.deb"
            else
                ensure_success $sh_c "cp ${aapv_tar} ./"
            fi
            # todo 
            # ensure_success $sh_c "dpkg -i apparmor_4.0.1-0ubuntu1_amd64.deb"
        fi
    fi

    # opy pre-installation dependency files 
    if [ -d /opt/deps ]; then
        ensure_success $sh_c "mv /opt/deps/* ${BASE_DIR}"
    fi
}

echo "---precheck_os---"
precheck_os
exit