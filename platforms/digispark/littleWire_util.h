#ifndef LITTLEWIRE_UTIL_H
#define LITTLEWIRE_UTIL_H

#ifdef _WIN32
  #include <windows.h>
#else
  #include <unistd.h>
#endif

/* Delay in miliseconds */
void delay(unsigned int duration);

#endif
