#ifndef CHECK_H
#define CHECK_H
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
int getLen(char *src,char separator);
void split(char *src, const char *separator, char **dest);
int Check(char *src, char *dst);
#endif /* CHECK_H */