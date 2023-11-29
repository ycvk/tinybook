-- 1, 2, 3, 4, 5, 6, 7 这是你的元素
-- ZREMRANGEBYSCORE key1 0 6
-- 7 执行完之后, 1, 2, 3, 4, 5, 6 这些元素就被删除了
-- 限流对象的键，通常是IP地址
local key = KEYS[1]
-- 窗口大小（单位：毫秒）
local window = tonumber(ARGV[1])
-- 阈值，窗口内允许的最大请求数
local threshold = tonumber(ARGV[2])
-- 当前时间（单位：毫秒）
local now = tonumber(ARGV[3])
-- 窗口的起始时间
local min = now - window

-- 移除旧的请求时间戳
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)

-- 获取当前窗口内的请求数量
local cnt = redis.call('ZCOUNT', key, min, now)

-- 如果请求数量没有超过阈值，则添加新的时间戳并设置键的过期时间
if cnt >= threshold then
    return 1
end

redis.call('ZADD', key, now, now)
-- 只有key是新的,才设置过期时间 因为如果key已经存在,那么它的过期时间每次都会被重置
if cnt == 0 then
    redis.call('PEXPIRE', key, window)
end

return 0
