# email-forward-parser
Parses forwarded emails and extracts original content.

This is a Go port of [crisp-oss/email-forward-parser](https://github.com/crisp-oss/email-forward-parser).

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
