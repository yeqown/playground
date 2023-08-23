LIBNAME = mylib.so
SRC     = mylib.c
LUA_INCLUDE_DIR = /opt/homebrew/Cellar/lua/5.4.6/include/lua5.4
LUA_LIB_DIR     = /opt/homebrew/Cellar/lua/5.4.6/lib

build:
	# compile the shared library with cc on macOS
	cc -bundle -undefined dynamic_lookup -o $(LIBNAME) $(SRC) -I$(LUA_INCLUDE_DIR) -L$(LUA_LIB_DIR) -llua