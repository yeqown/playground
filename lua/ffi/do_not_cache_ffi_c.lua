local ffi = require("ffi")

ffi.cdef[[
    int printf(const char *fmt, ...);
]]

-- 缓存 printf 函数
-- local cprintf = ffi.C.printf
-- for i = 1, 10000 do
-- local cint = ffi.new("int", i)
--     cprintf("Hello %d\n", cint)
-- end

-- 不缓存 printf 函数
for i = 1, 10000 do
    ffi.C.printf("Hello %d\n", ffi.new("int", i))
end

-- 缓存 ffi.C
-- local C = ffi.C
-- for i = 1, 10000 do
--     C.printf("Hello %d\n", ffi.new("int", i))
-- end

-- cache ffi.C.pringf:
-- Cost: 0:00.05 Secs
-- Memory: 1424 KB
-- 
-- do not cache ffi.C.pringf:
-- Cost: 0:00.01 Secs
-- Memory: 1424 KB