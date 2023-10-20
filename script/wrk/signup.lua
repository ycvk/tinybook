wrk.method = "POST"
wrk.headers["Content-Type"] = "application/json"

local thread_id = 0

local f = assert(io.open("/dev/urandom", "r"))
local d = f:read(4)
f:close()

function setup(thread)
    thread:set("id", thread_id)
    thread_id = thread_id + 1
end

function init(args)
    math.randomseed(os.time() + d:byte(1) + (d:byte(2) * 256) + (d:byte(3) * 65536) + (d:byte(4) * 16777216))
end

-- 生成uuid 使用随机数 但是这个不是v4版本
function create_uuid()
    local template = 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
    return string.gsub(template, 'x', function(c)
        local v = (c == 'x') and math.random(0, 0xf) or math.random(8, 0xb)
        return string.format('%x', v)
    end)
end

-- 按位与
function rshift(x, n)
    return math.floor(x / 2 ^ n)
end

-- 按位或
function band(x, y)
    local z, i = 0, 1
    while x > 0 and y > 0 do
        if x % 2 == 1 and y % 2 == 1 then
            z = z + i
        end
        x, y, i = rshift(x, 1), rshift(y, 1), i * 2
    end
    return z
end

-- 生成uuid v4版本 由于lua5.1不支持位运算，所以只能手动实现
function uuid()
    local template = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function(c)
        local t = band(os.time(), 0xFFFFFFF) + math.random()
        local v = c == 'x' and band(rshift(t * 16, math.random(0, 12)), 15) or band(rshift(t * 64, math.random(0, 10)), 15)
        return string.format('%x', v)
    end)
end

function request()
    local thread_specific_id = id  -- 假设 "id"已由 setup 函数设置
    local email = string.format("%s%d@qq.com", create_uuid(), thread_specific_id)
    local body = string.format('{"email":"%s", "password":"hello#world123", "confirmPassword": "hello#world123"}', email)
    return wrk.format('POST', wrk.path, wrk.headers, body)
end

function response(status, headers, body)
    -- Add any logic here if needed
end
