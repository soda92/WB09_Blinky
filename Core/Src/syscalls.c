/**
 * @file syscalls.c
 * @brief Minimal implementation of syscalls to silence Newlib linker warnings.
 *
 * This file provides weak dummy implementations for the system calls expected
 * by Newlib-Nano. They do nothing and return generic error codes or success
 * indicators as appropriate for a bare-metal environment.
 */

#include <sys/stat.h>
#include <unistd.h>
#include <errno.h>

/* Undefine errno to avoid conflicts if it's a macro */
#undef errno
extern int errno;

/* Environment variable pointer */
char *__env[1] = { 0 };
char **environ = __env;

/*
 * _close
 */
__attribute__((weak)) int _close(int file)
{
    return -1;
}

/*
 * _fstat
 */
__attribute__((weak)) int _fstat(int file, struct stat *st)
{
    st->st_mode = S_IFCHR;
    return 0;
}

/*
 * _isatty
 */
__attribute__((weak)) int _isatty(int file)
{
    return 1;
}

/*
 * _lseek
 */
__attribute__((weak)) int _lseek(int file, int ptr, int dir)
{
    return 0;
}

/*
 * _read
 */
__attribute__((weak)) int _read(int file, char *ptr, int len)
{
    return 0;
}

/*
 * _exit
 * Note: This might already be provided by Newlib, but sometimes it's missing.
 */
__attribute__((weak)) void _exit(int status)
{
    while (1);
}

/*
 * _kill
 */
__attribute__((weak)) int _kill(int pid, int sig)
{
    errno = EINVAL;
    return -1;
}

/*
 * _getpid
 */
__attribute__((weak)) int _getpid(void)
{
    return 1;
}
