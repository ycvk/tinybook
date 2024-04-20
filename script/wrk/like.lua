wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"
-- 记得修改这个，你在登录页面登录一下，然后复制一个过来这里
wrk.headers["Authorization"] = "Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDI1NjY1OTIsInVpZCI6MSwic3NpZCI6ImQ0YmNlMGIyLTQ4NGYtNDVhYi05ZWFhLWY4ZmU5M2QzMmI1YyIsInVzZXJBZ2VudCI6Ik1vemlsbGEvNS4wIChNYWNpbnRvc2g7IEludGVsIE1hYyBPUyBYIDEwXzE1XzcpIEFwcGxlV2ViS2l0LzUzNy4zNiAoS0hUTUwsIGxpa2UgR2Vja28pIENocm9tZS8xMTkuMC4wLjAgU2FmYXJpLzUzNy4zNiJ9.j-UBwnQ8-cEWiJ_MTq56sDYhFHzku4zHgTyzSySgjLrSXOMl6LeDHg0KHvTR2fCkxq0Y0wYqOysiR9_Bgn3ZfQ"
wrk.body = '{"id": 1228636373573963776, "like": true}'