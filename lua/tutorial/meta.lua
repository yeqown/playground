-- meta-like examples

-- t represents a table
local t = { a = 1 }

-- meta_table represents a meta table, which is also a table
-- meta methods pre-defined in lua, reference: https://www.lua.org/manual/5.3/manual.html#2.4
local meta_table_of_t = {
    -- __add will be called when t + 1 is called
    __add = function(a, b)
        return a.a + b
    end,

    -- __index will be called when t[k] is called and t[k] is nil
    __index = function(t, k)
        return k
    end,

    -- __newindex will be called when t[k] = v
    __newindex = function(t, k, v)
        print("new index", k, v)
    end
}

-- setmetatable(t, meta_table_of_t) set meta_table_of_t as the meta table of t
setmetatable(t, meta_table_of_t)

-- print(t + 1) -- 2
print(t + 1) -- 2

-- print(t.b) -- b
print(t.b) -- b


------------------------------------ metatable work as class ------------------------------------
local bag = {}
local bag_meta_table = {
    getname = function(self)
        return self.name
    end,
}
bag_meta_table["__index"] = bag_meta_table

function bag.new(self, bag_name)
    local o = {name= bag_name or "null"}
    setmetatable(o, bag_meta_table)
    return o
end

local b = bag:new("bag")
print("bagname is:defined", b:getname()) -- bag
