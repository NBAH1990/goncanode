# goncanode

Simplify https://github.com/malikzh/NCANode usage for Go.

Usage:
```go
v := types.NCAnodeV30                   // or types.NCAnodeV10
nH := goncanode.Create(entities.Options{
    ServiceUrl: conf.NcaNode.ServiceUrl,// http://127.0.0.1:14579
    P12base64: conf.NcaNode.P12Base64,  // base64 encoded p12 cert
    P12pass:   conf.NcaNode.P12Pass,    // p12 cert password
    Timeout: 1500 * time.Millisecond,   // context waiting timeout
    Version: &v,                        // NCANode version (differences in API)
})

sr, err := nH.SignWithSecurityHeader(r.Context(), xmlString, types.GOST34311)
```
