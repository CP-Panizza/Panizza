#include "check.h"

void split(char *src, const char *separator, char **dest)
{
    /*
        src 源字符串的首地址(buf的地址)
        separator 指定的分割字符
        dest 接收子字符串的数组
        num 分割后子字符串的个数
    */
    char *pNext;
    pNext = (char *)strtok(src,separator); //必须使用(char *)进行强制类型转换(虽然不写有的编译器中不会出现指针错误)
    while(pNext != NULL) {
        *dest++ = pNext;
        pNext = (char *)strtok(NULL,separator);  //必须使用(char *)进行强制类型转换
    }
}

int getLen(char *src,char separator){
    char *strp = src;
    int count = 0; //count '/' in string
    while(*strp != '\0'){
        if(*strp == separator) count++;
        strp++;
    }
    return count;
}

//检测url和路由是否符合规则
int Check(char *src, char *dst){
    int len_src = getLen(src, '/');
    int len_dst = getLen(dst, '/');
    if(len_src != len_dst) return 0;
    char **src_list = (char**)malloc(sizeof(char*) * len_src);
    split(src, "/", src_list);
    char **dst_list = (char**)malloc(sizeof(char*) * len_dst);
    split(dst, "/", dst_list);
    int i = 0;
    for (; i < len_src;i++) {
        if(strstr(src_list[i], ":")!= NULL){
            continue;
        }

        if(strcmp(src_list[i], dst_list[i])){
           return 0;
        }
    }
    free(src_list);
    free(dst_list);
    return 1;
}

