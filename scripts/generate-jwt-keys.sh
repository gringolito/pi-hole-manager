#!/bin/bash

# rsa-256
openssl genrsa -out rsa-256-private.pem 2048
openssl rsa -in rsa-256-private.pem -pubout -outform PEM -out rsa-256-public.pem
openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in rsa-256-private.pem -out rsa-256-private-pkcs8.pem

# ecdsa-256
openssl ecparam -name prime256v1 -genkey -noout -out ecdsa-256-private.pem
openssl ec -in ecdsa-256-private.pem -pubout > ecdsa-256-public.pem
openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in ecdsa-256-private.pem -out ecdsa-256-private-pkcs8.pem

# ecdsa-384
openssl ecparam -name secp384r1 -genkey -noout -out ecdsa-384-private.pem
openssl ec -in ecdsa-384-private.pem -pubout > ecdsa-384-public.pem
openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in ecdsa-384-private.pem -out ecdsa-384-private-pkcs8.pem

# ecdsa-512
openssl ecparam -name secp521r1 -genkey -noout -out ecdsa-512-private.pem
openssl ec -in ecdsa-512-private.pem -pubout > ecdsa-512-public.pem
openssl pkcs8 -topk8 -inform PEM -outform PEM -nocrypt -in ecdsa-512-private.pem -out ecdsa-512-private-pkcs8.pem
