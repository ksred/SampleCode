# Visa Samples - Go

The examples illustrate usage of the Visa Developer API using Go. Currently only VisaDirect has been completed, more samples will be added soon.

## Usage

- Log on to https://developer.visa.com/, go to the Dashboard and click on your app name
- Copy the User ID and Password from the Keys/APIs tab to any text editor
- Download the cert.pem from app details on VDP portal (should be visible under the Certificates when you click on the app name in Dashboard and go to Keys/APIs)
- Edit `SetUserPassword("user_id", "user_password")` in `VisaDirectSample_test` with the User ID and Password from step 2
- The default certificate paths are:
```
    SSL_PUBLIC_KEY_PATH = "certs/application.crt"
    SSL_PRIVATE_KEY_PATH = "certs/application.pem"
    SSL_CAPRIVATE_KEY_PATH = "certs/VDPCA-SBX.pem"
```

  - If defaults are different from above, edit `SetCertPaths("pubKeyPath", "pvtKeyPath", "caPvtKeyPath")` in `VisaDirectSample_test` with the certificate paths.

- To run VisaDirect, go to the folder containing VisaDirectSample.go and run the command:

`go test`

You should see response from the VisaDirect API calls as the tests complete.

To know more about generation of private key (exampled-key.pem) and CSR upload to create app, go to https://developer.visa.com/vdpguide#gettingStarted

## Visa package
There is a Visa package for Go [currently in development](https://github.com/ksred/visa).

## Issues
If there are any issues, please create an issue and tag @ksred.
