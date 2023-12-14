wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDIzMjM5NTYsInVpZCI6MSwic3NpZCI6ImRiODhhZjc1LTU5YzAtNDJhZi05ZmU4LWM0NzA1ZDBhYjI2MSIsInVzZXJBZ2VudCI6Ik1vemlsbGEvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMTkuMC4wLjAgU2FmYXJpLzUzNy4zNiJ9.LUOENqmsSZXj9bCC7Jg6AaY7D4B0im8SiyNlJlONLJ30SBg0O8Y-ULp5M7xhG1ap4Iq_mjrCe4a7qZm5tavXjw"