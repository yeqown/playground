local ffi = require("ffi")

ffi.cdef[[
    void Sleep(int ms);
    int poll(struct pollfd *fds, unsigned long nfds, int timeout);
]]

print("ffi.os: " .. ffi.os)
local sleep
if ffi.os == "Windows" then
    function sleep(s)
        ffi.C.Sleep(s * 1000)
    end
else
    function sleep(s)
        ffi.C.poll(nil, 0, s * 1000)
    end
end

for i = 1, 10 do
    io.write("."); io.flush()
    sleep(0.09)
end
io.write("\n")