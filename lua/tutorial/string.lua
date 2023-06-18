-- ABCD1234!!
-- 0x41 0x42 0x43 0x44 0x31 0x32 0x33 0x34 0x21 0x21
-- any byte can be stored in a string, even '\0' 0x00
-- All strings contains two index fields: 
-- 1. a positive integer index, which is the position of the character in the string
-- 2. a negative integer index, which is the position of the character counted from the end of the string

local s = "ABCD1234!!"
print(s[1], s[-#s]) -- A !
print(#s) -- 10
print(s[1] == s[-#s]) -- true

print(string.byte(s, 1)) -- 65 == 0x41, it'stored 'A'
-- strings.byte(s, i, j) i represent the start index, j represent the end index, if j is not given, it will be the same as i
-- if j == -1, it will be the same as #s
print(s:byte(1, -1)) -- 65 66 67 68 49 50 51 52 33 33