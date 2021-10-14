#!/bin/sh

export OPENSSL_CONF="../configs/openssl.conf"

certs_dir="certs"

if [ ! -d "$certs_dir" ]
then
    echo "Creating certs directory"
    mkdir "$certs_dir"
fi

cd "$certs_dir"

if [ ! -d "private" ]
then
    mkdir "private"
fi

error () {
    echo "Unable to create certificates" >&2
    exit 1
}

create_ca () {
    CN="AMIGA7 CA"; export CN
    DN_SECTION="dn_ca"; export DN_SECTION

    openssl req -new -x509 -nodes \
        -out tsa_ca.pem -keyout private/tsa_ca_key.pem
    test $? != 0 && error
}

create_tsa_cert () {
    EXT=$2
    CN=$1; export CN
    DN_SECTION="dn_cert"; export DN_SECTION

    openssl req -new \
        -out tsa_req.pem -keyout private/tsa_key.pem
    test $? != 0 && error

    echo Using extension $EXT
    openssl x509 -req \
        -in tsa_req.pem -out tsa_cert.pem \
        -CA tsa_ca.pem -CAkey private/tsa_ca_key.pem -CAcreateserial \
        -extfile $OPENSSL_CONF -extensions $EXT
    test $? != 0 && error
}

echo "Creating CA"
create_ca

echo "Creating TSA server cert"
create_tsa_cert "TSA CERT" tsa_cert

exit 0
