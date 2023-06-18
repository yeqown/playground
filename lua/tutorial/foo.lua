require('bar')
require('bar')
require('bar')
require('bar')
require('bar') -- prints "I'm bar.lua" only once

local r = require('package.bar2')
print(r)

-- package.path is a global variable, it's a string that contains the path of lua modules
-- It includes the current directory, and the directories in the environment variable LUA_PATH
-- "./?.lua" represents the current directory
print(package.path)

-- To run this file, type:
-- lua bar.lua