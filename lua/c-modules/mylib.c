#include <lauxlib.h>
#include <lua.h>
#include <stdio.h>

#ifndef VERSION
#define VERSION "1.0.0"
#endif

#ifndef NAME
#define NAME "mylib"
#endif

int add(lua_State *L) {
  int a = luaL_checknumber(L, 1);
  int b = luaL_checknumber(L, 2);
  lua_pushnumber(L, a + b);
  return 1;
}

int luaopen_mylib(lua_State *L) {
  luaL_Reg lib[] = {{"add", add}, {NULL, NULL}};
  luaL_newlib(L, lib);

  // 设置全局变量
  lua_pushvalue(L, -1);
  lua_setglobal(L, "mylib");

  lua_pushliteral(L, VERSION);
  lua_setfield(L, -2, "_VERSION");
  lua_pushliteral(L, NAME);
  lua_setfield(L, -2, "_NAME");

  printf("Hello, world you have import mylib.c\n");
  return 1;
}