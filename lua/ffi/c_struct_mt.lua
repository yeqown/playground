local ffi = require("ffi")

ffi.cdef[[
    typedef struct { double x, y; } point_t;
]]

local point
local mt = {
    __add = function(a, b)
        return point(a.x + b.x, a.y + b.y)
    end,

    __len = function(a)
        return math.sqrt(a.x * a.x + a.y * a.y)
    end,

    __index = {
        area = function(a)
            return a.x * a.y
        end
    },
}

-- 参见 ffi.metatype 的文档, 绑定元表到 point_t 类型
point = ffi.metatype("point_t", mt)

local a = point(3, 4)
print(a.x, a.y, #a, a:area())
-- 3 4 5 12
local b = a + point(0.5, 8)
print(#b)
-- 12.5