#!/bin/bash
declare -r bin_name=$(basename "${0}")

function print_help() {
	echo -e "Usage: ${bin_name} [-h] [-p PREFIX] KEY_TYPE"
	echo -e ""
    echo -e "Available key-types: rsa-256, rsa-384, rsa-512, ecdsa-256, ecdsa-384, ecdsa-512"
	echo -e "Options:"
    echo -e "  -p PREFIX Set output filename prefix (Defaults to: KEY_TYPE)"
	echo -e "  -h        Print this help message"
	exit "${1}"
}

function verify_command() {
    local cmd=${1}
    if ! command -v ${cmd} > /dev/null 2>&1; then
        echo >&2 "${cmd} is required but not installed. Please install ${cmd} and try again."
        exit 1
    fi
}

function generate_rsa() {
    yes | ssh-keygen -t rsa ${rsa_extra_params} -m PEM -b 4096 -N "" -f ${private_filename} > /dev/null
    rm -f "${private_filename}.pub"
    openssl rsa -in "${private_filename}" -pubout -outform PEM -out "${public_filename}" 2> /dev/null
}

function generate_ecdsa() {
    local random=$(openssl rand -hex 8)
    local certificate="ecdsa.${random}.pem"
    openssl ecparam ${ecdsa_extra_params} -genkey -noout -out "${certificate}"
    openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in "${certificate}" -out "${private_filename}"
    openssl ec -in "${certificate}" -pubout > "${public_filename}" 2> /dev/null
    rm -f "${certificate}"
}

verify_command openssl
verify_command ssh-keygen

while getopts p:h OPT; do
	case "${OPT}" in
    p) prefix="${OPTARG}" ;;
	h) print_help 0 ;;
    *) print_help 1 ;;
	esac
done

key_type="${@:$OPTIND:1}"
test -z "${key_type}" && print_help 1

known_types=("rsa-256" "rsa-384" "rsa-512" "ecdsa-256" "ecdsa-384" "ecdsa-512")
if [[ ! " ${known_types[*]} " =~ " ${key_type} " ]]; then
    echo -e "Invalid key type: ${key_type}"
    print_help 1
fi

test -z "${prefix}" && prefix="jwt-${key_type}"
private_filename="${prefix}.pem"
public_filename="${prefix}.pub"

algorithm=$(echo ${key_type} | cut -d '-' -f 1)
shalen=$(echo ${key_type} | cut -d '-' -f 2)

if [[ "${shalen}" -eq 256 ]]; then
    rsa_extra_params="-E sha256"
    ecdsa_extra_params="-name prime256v1"
elif [[ "${shalen}" -eq 384 ]]; then
    rsa_extra_params="-E sha384"
    ecdsa_extra_params="-name secp384r1"
else
    rsa_extra_params="-E sha512"
    ecdsa_extra_params="-name secp521r1"
fi

# set -x
set -e

if [[ "${algorithm}" == "rsa" ]]; then
    generate_rsa
else
    generate_ecdsa
fi

echo -e "JWT signing keys sucessfully generated"
echo -e "Private key:\t${private_filename}"
echo -e "Public key:\t${public_filename}"
