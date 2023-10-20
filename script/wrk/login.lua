wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"
-- 这个要改为你的注册的数据
wrk.body = '{"email":"4cc9ad7c4e9fb7f5a9ae9371ef856bff0@qq.com", "password": "hello#world123"}'