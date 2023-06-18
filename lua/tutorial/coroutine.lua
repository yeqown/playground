-- coroutine module tutorial
-- refence: http://www.lua.org/

local routine = coroutine.create(function()
    print("hello")

    -- yield the coroutine, which means the coroutine is suspended
    -- and the coroutine will be resumed when the coroutine.resume is called
    -- 1,2,3 are the return values of coroutine.resume
    -- r2 is the parameter of coroutine.resume
    local r2 = coroutine.yield(1,2,3)
    print("r2: ", r2) -- 4
end)

print(type(routine)) -- thread
-- coroutine.status(routine) -- get the status of coroutine
print(coroutine.status(routine)) -- suspended

-- yet, the coroutine is not running, so we need to resume it
-- now, resume the coroutine

print("coroutine.resume: ", coroutine.resume(routine)) -- true 1 2 3
-- hello
-- coroutine.resume:  true 1 2 3

coroutine.resume(routine, 4) -- r2:  4
