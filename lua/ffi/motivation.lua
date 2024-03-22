-- 1. Load ffi module
local ffi = require "ffi"

-- 2. Add a C declaration for printf
ffi.cdef[[
    int printf(const char *fmt, ...);
]]

-- 3. Call the named C function printf
ffi.C.printf("Hello %s!\n", "world")

-- 4. run this file with luajit
-- luajit motivation.lua