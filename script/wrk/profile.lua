wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "RapidAPI/4.2.0 (Macintosh; OS X/14.0.0) GCDHTTPRequest"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc4MjUxODYsInVpZCI6NzUxLCJ1c2VyQWdlbnQiOiJSYXBpZEFQSS80LjIuMCAoTWFjaW50b3NoOyBPUyBYLzE0LjAuMCkgR0NESFRUUFJlcXVlc3QifQ.Ohjb8mvUtcIPsNAn_cNSZYNAuhRipfvW-rRRPDJqgzTJst_-rlxGStkZ0Za8lkyFaN8WCEgI_D-dMJHQWloyeA"