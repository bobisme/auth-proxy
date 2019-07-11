Authentication required

```
⟩ http :8080
HTTP/1.1 401 Unauthorized
Content-Length: 12
Content-Type: text/plain; charset=utf-8
Date: Thu, 11 Jul 2019 21:03:56 GMT
Server: fasthttp

Unauthorized
```

---

Dummy Authorization header

```
⟩ http :8080 'Authorization:Bearer x'
HTTP/1.1 200 OK
Content-Length: 84
Content-Type: text/plain; charset=utf-8
Date: Thu, 11 Jul 2019 21:04:13 GMT
Server: fasthttp

Hi there, bob@8o8.me!
user id = 123
client id =
scopes = []
client can not have fun
```

---

X-UserInfo stripped out

```
⟩ http :8080 'Authorization:Bearer x' 'X-Userinfo: {"uid":"123","sub":"iam@bobis.me","cid":"234","scopes":["a","b","fun"]}'
HTTP/1.1 200 OK
Content-Length: 84
Content-Type: text/plain; charset=utf-8
Date: Thu, 11 Jul 2019 21:14:13 GMT
Server: fasthttp

Hi there, bob@8o8.me!
user id = 123
client id =
scopes = []
client can not have fun
```

```
[00] {"level":"warn","uerInfoHeader":"{\"uid\":\"123\",\"sub\":\"iam@bobis.me\",\"cid\":\"234\",\"scopes\":[\"a\",\"b\",\"fun\"]}","time":"2019-07-11T17:14:13-04:00","message":"client tried to pass X-UserInfo header"}
```
