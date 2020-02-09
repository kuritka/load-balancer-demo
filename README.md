# load-balancer-demo

### Running locally 

install certificate by default go tool `generate_cert.go`
```bash
go install $GOROOT/src/crypto/tls/generate_cert.go
./generate_cert -host localhost
cp -f  key.pem /etc/lb/certs/key.pem
cp -f  key.pem /etc/lb/certs/cert.pem
```