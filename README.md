# DNS Check (dnsck)

dnsck performs DNS lookups and validates the connection with one
or more hosts. It returns the results in JSON format allowing
parsing with tools like jq (https://stedolan.github.io/jq/).

dnsck is written in Go (https://go.dev).

You can download executable versions from the releases area in GitHub:
https://github.com/andreburgaud/dnsck

## Usage

```
$ dnsck google.com
{
  "servers": [
    {
      "hostname": "google.com",
      "connection": {
        "address": "google.com:443",
        "success": true
      },
      "domain": {
        "name": "google.com",
        "ip_addresses": [
          "172.253.119.101",
          "172.253.119.102",
          "172.253.119.113",
          "172.253.119.139",
          "172.253.119.100",
          "172.253.119.138",
          "2607:f8b0:4001:c23::65",
          "2607:f8b0:4001:c23::71",
          "2607:f8b0:4001:c23::8a",
          "2607:f8b0:4001:c23::66"
        ],
        "canonical_name": "google.com.",
        "dns_text_records": [
          "facebook-domain-verification=22rm551cu4k0ab0bxsw536tlds4h95",
          "MS=E4A68B9AB2BB9670BCE15412F62916164C0B20BB",
          "v=spf1 include:_spf.google.com ~all",
          "docusign=05958488-4752-4ef2-95eb-aa7ba8a3bd0e",
          "globalsign-smime-dv=CDYX+XFHUw2wml6/Gb8+59BsH31KzUr6c1l2BPvqKX8=",
          "google-site-verification=wD8N7i1JTNTkezJ49swvWW48f8_9xveREV4oB-0Hf5o",
          "onetrust-domain-verification=de01ed21f2fa4d8781cbc3ffb89cf4ef",
          "webexdomainverification.8YX6G=6e6922db-e3e6-4a36-904e-a805c28087fa",
          "apple-domain-verification=30afIBcvSuDV2PLX",
          "google-site-verification=TV9-DBe4R80X4v0M4U_bd_J9cpOJM0nikft0jAgjmsQ",
          "docusign=1b0a6754-49b1-4db5-8540-d2c12664b289",
          "atlassian-domain-verification=5YjTmWmjI92ewqkx2oXmBaD60Td9zWon9r6eakvHX6B77zzkFQto8PQ9QsKnbf4I"
        ]
      }
    }
  ],
  "count": 1,
  "count_connection_error": 0,
  "count_dns_error": 0
}
```

## Development

To build `dnsck`, you can use the `justfile` available in the root folder of the project or execute:

```
$ go build .
```

The `justfile` requires to have the command runner `just`, https://github.com/casey/just, on your machine.

# License

MIT