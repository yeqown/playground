local core = require("apisix.core")
local plugin_name = "ip2location"

local schema = {
    type = "object",
    properties = {},
}

local _M = {
    version = 0.1,
    priority = 0,
    name = plugin_name,
    schema = schema,
}

function _M.check_schema(conf)
    return core.schema.check(schema, conf)
end

function _M.rewrite(conf, ctx)
    local header_a = core.request.header(ctx, "X-Real-IP")
    if header_a then
        core.request.set_header(ctx, "X-Real-IP2", header_a)
    end
end

return _M
