-- all nginx request will be redirected to this file
-- reply with a simple string "hello world"

-- read from mapping file
-- format: hostname redirect_hostname
-- example: baidu2.com baidu.com
-- local mapping = {}
-- local file = io.open("/etc/nginx/lua/mapping.config", "r")
-- for line in file:lines() do
--     local key, value = string.match(line, "(%S+)%s+(%S+)")
--     mapping[key] = value
-- end

-- if hostname is baidu2.com, then redirect to baidu.com
-- redirect to another url (302)
-- ngx.var.host is the hostname


if ngx.var.host == "baidu2.com" then
    ngx.redirect("http://baidu.com", 302)
else 
    -- if hostname is not baidu2.com, then reply with a string
    -- hostname and port
    ngx.say("hello world, hostname: ", ngx.var.host, " port: ", ngx.var.server_port)
end
