local ffi = require "ffi"

ffi.cdef[[
    typedef struct {
        uint8_t red;
        uint8_t green;
        uint8_t blue;
    } rgb_pixel_t;
]]

local function image_ramp_green(n)
    -- ffi.new equivalent to malloc
    local img = ffi.new("rgb_pixel_t[?]", n)
    local f = 255/(n-1)
    for i = 0, n-1 do
        img[i].red = 0
        img[i].green = i*f
        img[i].blue = 0
    end
    return img
end

local function image_to_grey(img, n)
    for i = 0, n-1 do
        local y = 0.3*img[i].red + 0.59*img[i].green + 0.11*img[i].blue
        img[i].red = y
        img[i].green = y
        img[i].blue = y
    end
end

local N = 400*400
local img = image_ramp_green(N)
for i = 1, 100 do
    image_to_grey(img, N)
end

-- The difference between data_structure.lua and c_data_structure.lua:
-- 
-- gtime -f "Cost: %E Secs\nMemory: %M KB" luajit data_structure.lua
-- gtime = gnu-time on Mac OS X, time on Linux
-- 
-- | Contrast Item | data_structure.lua          | c_data_structure.lua |
-- | :------------ | :-------------------------- | :------------------- |
-- | Cost          | (luajit)0.57s / (lua)11.72s | (luajit)0.03s        |
-- | Memory        | (luajit)30MB  / (lua)30MB   | (luajit)2MB          |
-- 
--  Conclusion:
-- 
--  1. luajit is much faster than lua
--  2. c data structure cost less memory than lua table
-- 
--  That's why we need to use ffi to access C data structure, 
--  in a word, it's faster and cost less memory than lua table.