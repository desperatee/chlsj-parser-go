# Charles session request header parser

Automatically convert request headers from a session file to Golang code

## Usage:
```
go run . -session=<filename>.chlsj

Example output:
2023/05/28 01:09:31 INFO Parsing request #0 url=https://warsawsneakerstore.com/order/payment
var parsed = http.Header{
        "Upgrade-Insecure-Requests": {"1"},
        "Sec-Ch-Ua-Platform": {"\"Windows\""},
        "Sec-Fetch-Mode": {"navigate"},
        "Sec-Fetch-Dest": {"empty"},
        "Referer": {"https://warsawsneakerstore.com/order/payment"},
        "Origin": {"https://warsawsneakerstore.com"},
        "Accept-Language": {"pl-PL,pl;q=0.9,en-US;q=0.8,en;q=0.7"},
        "Cookie": {""},
        "Sec-Ch-Ua-Mobile": {"?0"},
        "Content-Type": {"application/x-www-form-urlencoded"},
        "Accept": {"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"},
        "Sec-Fetch-Site": {"same-origin"},
        http.HeaderOrderKey: {"content-length","cache-control","sec-ch-ua","origin","sec-ch-ua-mobile","upgrade-insecure-requests","user-agent","content-type","accept","sec-ch-ua-platform","sec-fetch-site","sec-fetch-mode","sec-fetch-dest","referer","accept-encoding","accept-language","cookie"},
        http.PHeaderOrderKey: {":method",":authority",":scheme",":path"},
        "Cache-Control": {"max-age=0"},
        "Sec-Ch-Ua": {"\"Google Chrome\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\""},
        "User-Agent": {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"},
        "Accept-Encoding": {"gzip, deflate, br"},
}
```