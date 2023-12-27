-- 创建一个模拟的 request_handle 对象
request_handle = {}

function request_handle:new(cvn) -- 创建一个新的对象
    local o = {}
    setmetatable(o, self)
    self.__index = self
    o.headers = function()
        return {
            get = function(self, key)
                if key == ":path" then
                    return "/path?cvn=" .. cvn -- 返回模拟的请求头
                end
            end,
            add = function(self, key, value)
                print(key .. ": " .. value) -- 打印添加的头
            end
        }
    end
    return o
end

-- 测试的函数
function envoy_on_request(request_handle)
    local cvn = request_handle:headers():get(":path"):match("cvn=([%w%.%-]+)")
    if cvn == nil then
        cvn = "default"
    end
    request_handle:headers():add("x-client-cvn", cvn)
end

-- 调用测试的函数
local rh1 = request_handle:new("test")
envoy_on_request(rh1)
local rh2 = request_handle:new("v1.0.0")
envoy_on_request(rh2)
local rh3 = request_handle:new("v1.0.0-hotfix")
envoy_on_request(rh3)
local rh4 = request_handle:new("v1.0.0-rc.1")
envoy_on_request(rh4)
