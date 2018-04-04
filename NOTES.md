# Steps

## Obtain Service Key:

GET 
    https://olympus.itunes.apple.com/v1/app/config?hostname=itunesconnect.apple.com

Response:

```json
{
    "authServiceUrl": "https://idmsa.apple.com/appleauth",
    "authServiceKey": "..."
}
```

## Sign In

Signing in requires an email address and password. Two cookies will be returned
that you'll need to pass along in subsequent requests, `acn01` and `myacinfo`.

POST 
    https://idmsa.apple.com/appleauth/auth/signin

HEADERS:
    Content-Type: application/json
    X-Requested-With: XMLHttpRequest
    X-Apple-Widget-Key: [AUTH_SERVICE_KEY]
    Accept: application/json, text/javascript

BODY:
```json
{
    "accountName": "APPLE_ID",
    "password": "APPLE_PASSWORD",
    "rememberMe": true
}
```

Response Headers:
    
```
Set-Cookie: acn01=...
Set-Cookie: myacinfo=...
```

## Obtain Session

After signing in you need to obtain a session. The response has some relevant
team data but the important bit is the `itctx` cookie which you'll need in 
subsequent requests to convey your team context.

GET 
    https://olympus.itunes.apple.com/v1/session

COOKIES
     acn01=...
     myacinfo=...

Response Headers:

```
Set-Cookie: itctx=...
```

Response Body:

```json
{
    "user": {
        "fullName": "...",
        "emailAddress": "...",
        "prsId": "..."
    },
    "provider": {
        "providerId": 0,
        "name": "...",
        "contentTypes": ["SOFTWARE"]
    },
    "availableProviders": [ 
        {
            "providerId": 0,
            "name": "...",
            "contentTypes": [ "SOFTWARE" ]
        }, 
        {
            "providerId": 0,
            "name": "...",
            "contentTypes": [ "SOFTWARE" ]
        } 
    ]
}
```

## Switch Teams

This allows you to switch teams if you have more than one. The import bit here
is the `itctx` cookie returned. You'll want to use this for subsequent requests
so the correct team context is being conveyed. 

POST
    https://itunesconnect.apple.com/WebObjects/iTunesConnect.woa/ra/v1/session/webSession

COOKIES
     acn01=...
     myacinfo=...
     itctx=...

HEADERS:
    Content-Type: application/json; charset=utf-8
    
BODY:
```json
{
    "dsId": "USER_PRSID", 
    "contentProviderId": "PROVIDER_ID", 
    "ipAddress": null
}
```

Response Headers:

Set-Cookie: itctx=...

## Details

GET
    https://itunesconnect.apple.com/WebObjects/iTunesConnect.woa/ra/apps/manageyourapps/summary/v2

COOKIES
     acn01=...
     myacinfo=...
     itctx=...

## Applications

GET
    https://itunesconnect.apple.com/WebObjects/iTunesConnect.woa/ra/apps/manageyourapps/summary/v2?hostname=itunesconnect.apple.com
    
COOKIES
     acn01=...
     myacinfo=...
     itctx=...

BODY
```json
{
  "dsId": "...",
  "contentProviderId": "..."
}
```
