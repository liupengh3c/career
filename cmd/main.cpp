#include <iostream>
#include <thread>
#include <getopt.h>
 
int main(int argc, char *argv[]) 
{
    unsigned int n = std::thread::hardware_concurrency();
    std::cout << " concurrent threads are supported = " << n << std::endl;
    int opt;
    std::string str = "a::b:c:d";
    static struct option long_options[] =
    {  
        {"reqarg", required_argument,NULL, 'r'},
        {"optarg", optional_argument,NULL, 'o'},
        {"noarg",  no_argument,         NULL,'n'},
        {NULL,     0,                      NULL, 0},
    }; 
    int option_index;
    while((opt = getopt_long(argc,argv,str.data(),long_options,&option_index))!= -1)
    {  
        printf("opt = %c\t\t", opt);
        printf("optarg = %s\t\t",optarg);
        printf("optind = %d\t\t",optind);
        printf("argv[optind] =%s\t\t", argv[optind]);
        printf("option_index = %d\n",option_index);
    }  
}