local key = KEYS[1]
local cntKey = key .. ":cnt"
local val = ARGV[1]
local expireTime = tonumber(ARGV[2]) or 600
local minTtl = 60

local ttl = tonumber(redis.call("ttl", key))

-- -1 key没有设置过期时间
if ttl == -1 then
    return -1
elseif ttl < expireTime - minTtl then
    -- key设置的过期时间小于最小过期时间，说明已经可以再次发送验证码
    redis.call("setex", key, expireTime, val)
    redis.call("setex", cntKey, expireTime, 3)
    return 0
else
    return 1 -- 1表示还不能发送验证码
end
