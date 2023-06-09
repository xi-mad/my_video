# older ffmpeg doesn't have swscale. you might need to build a newer ffmpeg from source

CFLAGS=-W -Wall -O3 -s -D_FILE_OFFSET_BITS=64 -mtune=i686 -march=i386 -I/home/tuitfun/include -L/home/tuitfun/lib
LDFLAGS= -lavcodec -lavformat -lavcodec -lswscale -lavutil -lgd -lfreetype -ljpeg -lz -lm -lpthread
#LIBS=/home/tuitfun/lib/libavcodec.a /home/tuitfun/lib/libavformat.a /home/tuitfun/lib/libavutil.a /home/tuitfun/lib/libswscale.a
#LIBS=/usr/lib/libgd.a

all: mtn mtns

mtn: mtn.c Makefile
	$(CC) -o mtn mtn.c $(CFLAGS) $(LDFLAGS) $(LIBS)

mtns: mtn.c Makefile
	$(CC) -o mtns mtn.c -static $(CFLAGS) $(LDFLAGS) $(LIBS)

clean:
	rm -f mtn mtns
