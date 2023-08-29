-- 元类
Shape = { area = 0 }

-- 基础类方法 new
function Shape:new(o, side)
    o = o or {}
    setmetatable(o, self)
    self.__index = self
    side = side or 0
    self.area = side * side;
    return o
end

-- 基础类方法 printArea
function Shape:printArea()
    print("Shape 面积为 ", self.area)
end

-- 创建对象
local myshape = Shape:new(nil, 10)
myshape:printArea()

-- 继承
local Square = Shape:new()

-- 派生类方法 new
function Square:new(o, side)
    o = o or Shape:new(o, side)
    setmetatable(o, self)
    self.__index = self
    return o
end
-- 派生类方法 printArea
function Square:printArea()
    print("正方形面积为 ", self.area)
end

-- 创建对象 Square
local mysquare = Square:new(nil, 10)
mysquare:printArea()