cd certificates

read -p "Enter the client name: " client
read -p "Enter the client config file name: " clientconfig
read -p "Enter the CA certificate name: " rootca

openssl req -new -nodes -newkey rsa:4096 -keyout ./client-certificates/$client-key.pem -out ./client-certificates/$client.csr -config ./client-certificates/$clientconfig.cnf

openssl x509 -req -in ./client-certificates/$client.csr -copy_extensions=copy -CA ./root-ca-certificates/$rootca-cert.pem -CAkey ./root-ca-certificates/$rootca-key.pem -CAcreateserial -out ./client-certificates/$client-cert.pem -days 365 -sha256 
