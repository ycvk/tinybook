local key = KEYS[1]
local cntKey = key .. ":cnt"
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
if cnt == nil then
    -- 说明验证码已经过期了
    return -1
elseif cnt <= 0 then
    -- 说明验证码错误次数过多
    return -2
end

local code = redis.call("get", key)
if code == nil then
    -- 说明验证码已经过期了
    return -1
elseif code == expectedCode then
    -- 说明验证码正确
    redis.call("set", cntKey, 0)
    return 0
else
    -- 说明验证码错误
    redis.call("decr", cntKey)
    return 1
end
