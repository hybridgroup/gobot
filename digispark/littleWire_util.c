
#include <littleWire_util.h>

/* Delay in miliseconds */
void delay(unsigned int duration)
{
	#ifdef _WIN32
    Sleep(duration);
	#else
    usleep(duration*1000);
	#endif
}
