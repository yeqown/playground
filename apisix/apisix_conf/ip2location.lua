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
    -- "cloudfront-viewer-city": "Hong Kong" 设置到 X-City
    -- "cloudfront-viewer-country": "SG" 设置到 X-Country-Short Country-Short
    -- "cloudfront-viewer-country-name": "Singapore" 设置到 X-Country-Long
    -- "cloudfront-viewer-region-name": "Singapore" 设置到 X-Region

    local headers = ctx.var.headers
    local city = headers["cloudfront-viewer-city"]
    local country = headers["cloudfront-viewer-country"]
    local country_name = headers["cloudfront-viewer-country-name"]
    local region = headers["cloudfront-viewer-region-name"]

    if city then
        core.request.set_header("X-City", city)
    end

    if country then
        core.request.set_header("X-Country-Short", country)
        core.request.set_header("Country-Short", country)
    end

    if country_name then
        core.request.set_header("X-Country-Long", country_name)
    end

    if region then
        core.request.set_header("X-Region", region)
    end
end

return _M
