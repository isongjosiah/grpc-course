#!/bin/bash
#inspired from: https://github.com/grpc/grpc-java/tree/master/examples#generating-selfsigne-certificates... to be continued

#Output files
#ca.key: Certificate Authority private key file (this shouldn't be shared in real-life)
#ca.cert: Certificate Authority trust certificate (this should be share dwith users in real-life)
#server.key: Server private key, password protected (this shouldn't be shared)
#server.csr: Server certificate signing request(this should be shared with the CA owner)
#server.crt: Server certificate signed by the CA (this would be sent back by the CA owner) - Keep on server
#server.pem: Conversion of server.key into a format gRPC likes (this shouldn't be shared)

#Summary
# private files: ca.key, server.key, server.pem, server.crt
#"share" files: ca.crt (needed by the client), server.csr (needed by the CA)

# changes these CN's to match your host in your environment if needed.
SERVER_CN=localhost

#step 1; Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

#step 2: Generate the Server Private Key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

#step 3: Get a certificate signing request from the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

#step 4: Sign the certificate with the CA we created (It's called self signing) - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

#step 5: Convert the server certificate to .pem format (server.pem) usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem