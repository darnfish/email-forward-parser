# email-forward-parser
[![Go Reference](https://pkg.go.dev/badge/github.com/darnfish/email-forward-parser.svg)](https://pkg.go.dev/github.com/darnfish/email-forward-parser)
[![Build and Test](https://github.com/darnfish/email-forward-parser/actions/workflows/test.yml/badge.svg)](https://github.com/darnfish/email-forward-parser/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/darnfish/email-forward-parser/branch/main/graph/badge.svg?token=P8KQD92JZH)](https://codecov.io/gh/darnfish/email-forward-parser)

Parses forwarded emails and extracts original content.

This is a Go port of [crisp-oss/email-forward-parser](https://github.com/crisp-oss/email-forward-parser).

## Who uses it?

<table>
<tr>
<td align="center"><a href="https://pickupapp.io/"><img src="https://pickup.s3.darn.cloud/icons/256x256.png" height="64" /></a></td>
</tr>
<tr>
<td align="center">Pickup</td>
</tr>
</table>

## Installation
```
go get "https://github.com/darnfish/email-forward-parser"
```

## Usage
```go
import efp "github.com/darnfish/email-forward-parser"

result := efp.Read(body, subject)

log.Println(result.Forwarded) // true

log.Println(result)
// {
//   Forwarded: true,
//   Message: "Praesent suscipit egestas hendrerit.",
//   Email: {
//     Body: "Aenean quis diam urna.",
//     From: {
//       Name: "John Doe",
//       Address: "john.doe@acme.com"
//     },
//     To: [{
//       Name: "Bessie Berry",
//       Address: "bessie.berry@acme.com"
//     }],
//     CC: [{
//       Name: "Walter Sheltan",
//       Address: "walter.sheltan@acme.com"
//     }],
//     Subject: "Integer consequat non purus",
//     Date: "25 October 2021 at 11:17:21 EEST"
//   }
// }
```

## Licence
MIT
