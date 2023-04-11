#include <stdio.h>
#include <stdlib.h>

int main() {
    uint8_t *data = malloc(sizeof(uint8_t) * 4);
    *(uint32_t *)data = 0xabcd;
    printf("data: %s, data: %d \n", data, *(uint32_t*)data, data);
    return 0;
}
