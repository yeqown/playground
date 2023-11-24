#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <sys/inotify.h>

#define EVENT_SIZE  (sizeof (struct inotify_event))
#define BUF_LEN     (1024 * (EVENT_SIZE + 16))

int main(int argc, char *argv[]) {
    int fd, wd;
    char buf[BUF_LEN];
    ssize_t len;
    struct inotify_event *event;

    fd = inotify_init();
    if (fd == -1) {
        perror("inotify_init");
        exit(EXIT_FAILURE);
    }

    wd = inotify_add_watch(fd, "/tmp/test", IN_MODIFY);
    if (wd == -1) {
        perror("inotify_add_watch");
        exit(EXIT_FAILURE);
    }

    printf("Watching /tmp/test for changes...\n");

    while (1) {
        len = read(fd, buf, BUF_LEN);
        if (len == -1 && errno != EAGAIN) {
            perror("read");
            exit(EXIT_FAILURE);
        }

        if (len <= 0) {
            continue;
        }

        event = (struct inotify_event *) buf;
        if (event->mask & IN_MODIFY) {
            printf("File /tmp/test was modified!\n");
        }
    }

    inotify_rm_watch(fd, wd);
    close(fd);

    return 0;
}
