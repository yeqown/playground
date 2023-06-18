#include<stdio.h>
#include<errno.h>
#include<unistd.h>

#define MAX_WRITE_COUNT 100

int main() {
    // write to a pipe file (./waf.pipe)
    FILE *fp = fopen("./waf.pipe", "w");
    int unsigned count = 0;
    while (1)
    {
        // write to pipe file
        fprintf(fp, "hello world\n");

        // detect error
        if (ferror(fp))
        {
            printf("error %d\n", errno);
            break;
        }

        // flush to disk
        fflush(fp);
        count++;
        if (count >= MAX_WRITE_COUNT) {
            break;
        }

        // delay 1s
        sleep(1);
    }

    fclose(fp);
    // remove file
    if (remove("./waf.pipe") != 0) {
        printf("remove file failed\n");
    }


    return 0;
}