#include <stdio.h>
#include <errno.h>
#include <unistd.h>

int main()
{
    // read from a pipe file (./waf.pipe), split by '\n' and loop print each line
    FILE *fp = fopen("./waf.pipe", "r");

    // TODO: read forever until EOF
    while (1)
    {
        char buf[1024];
        if (fgets(buf, 1024, fp) != NULL)
        {
            printf("%s", buf);
        }

        // detect EOF
        if (feof(fp))
        {
            printf("read EOF, quit now\n");
            break;
        }

        // detect error
        if (ferror(fp))
        {
            printf("error %d\n", errno);
            break;
        }

        // delay 1s
        sleep(1);
    }

    fclose(fp);

    return 0;
}

// compile: gcc read.c -o read
// run: ./read