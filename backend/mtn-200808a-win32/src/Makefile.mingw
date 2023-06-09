CC=gcc
CFLAGS=-W -Wall -O3 -s -D_FILE_OFFSET_BITS=64 -mms-bitfields -mtune=i686 -march=i386 -I/usr/local/include -L/usr/local/lib
LDFLAGS= -lavcodec -lavformat -lavcodec -lswscale -lavutil -lz
#LDFLAGS= -lz
LIBS=/usr/local/lib/bgd.lib
#LIBS=/d/gd-2.0.34-win32/lib/bgd.lib /usr/local/bin/avcodec-51.dll /usr/local/bin/avformat-51.dll /usr/local/bin/avutil-49.dll /usr/local/bin/swscale-0.dll
#LIBS=/usr/local/lib/bgd.lib /usr/local/lib/avcodec-51.lib /usr/local/lib/avformat-52.lib /usr/local/lib/avutil-49.lib /usr/local/lib/swscale-0.lib

all: mtn

mtn: mtn.c Makefile
	$(CC) -o mtn mtn.c $(CFLAGS) $(LDFLAGS) $(LIBS)

mtns: mtn.c Makefile
	$(CC) -o mtns mtn.c -static $(CFLAGS) $(LDFLAGS) $(LIBS)

test: test.c Makefile
	$(CC) -o test test.c $(CFLAGS) $(LDFLAGS) $(LIBS)

clean:
	rm -f mtn.exe mtns.exe test.exe

