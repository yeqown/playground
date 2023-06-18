-- 基本类型演示
-- 1. nil
print(type(nil))
local c
print("c", c) -- nil

-- 2. boolean
print(type(true))
print(type(false))
print(type(0))
print(type(1))
print(1~=2)
print(1==2)
print(true and false)
print(true or false)
print(not 1)
print(not nil)
print(not false)
print(not 0)
print(1 >= 2)

-- 3. number
print(1.0)
print(1e10)
print(1)
print(0xff) -- 255

-- 4. string
print("hello world")
print('hello' .. ' world')
print('hello' .. 1)
print('hello' == 'hello')
print('hello' == 'world')

s = string.char(97, 98, 99) -- abc
print(s, #s) -- abc 3
s = string.byte(s, 1) -- 97

-- 5. table
local t = {0, 0} -- as array, index start from 1
t[1] = 1
t[2] = 2
print(t[1], t[2]) -- 1 2
print(#t) -- 2
table.insert(t, 3)
table.insert(t, 4, 1)
print(#t) -- 4
print("--------------------")
for i, v in ipairs(t) do
    print(i, v)
end
print("--------------------")
for i=1, #t do
    print(i, t[i])
end

d = {a = 1, b = 2} -- as map, index can be any type
print(d.a, d.b) -- 1 2
print(d['a'], d['b']) -- 1 2
for k, v in pairs(d) do
    print(k, v)
end
print(#d) -- 2

-- _G is a global table
print(_G)

-- 6. function
local function print_hello()
    print("hello world")
end

local add = function(a, b)
    print("add", a, "+", b, "=", a + b)
    return a + b
end

local function return_multiple_values()
    return 1, 2, 3
end

print_hello()
print(add(1, 2))
print(return_multiple_values())

-- 6.1 for loop
for i = 1, 10 do
    print(i)
end

for i = 1, 10, 2 do -- start from 1, step 2, end with 10
    print(i)
    i = i + 1 -- no effect
end
-- 6.2 if branch
if 1 == 1 then
    print("1 == 1")
elseif 1 == 2 then
    print("1 == 2")
else
    print("1 != 1 && 1 != 2")
end
-- while loop
local i = 1
while i <= 10 do
    if i == 5 then
        break
    end
    print(i)
    i = i + 1
end
