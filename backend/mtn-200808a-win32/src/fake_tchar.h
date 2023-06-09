/* 
 * fake_tchar.h
 *
 * Modified from mingw's tchar.h. This file should be used only in systems 
 * without the real tchar.h. This file maps _t* functions to the normal ones.
 *
 * This file has no copyright assigned and is placed in the Public Domain.
 * This file is a part of the mingw-runtime package.
 * No warranty is given; refer to the file DISCLAIMER within the package.
 *
 */

#ifndef FAKE_TCHAR_H
#define FAKE_TCHAR_H

typedef char    TCHAR;

#define _TEOF EOF

#define _TEXT(q)    q
#define _T(q)       q

/*
 * Non-unicode (standard) functions
 */

#define _tprintf    printf
#define _ftprintf   fprintf
#define _stprintf   sprintf
#define _sntprintf  snprintf
#define _vtprintf   vprintf
#define _vftprintf  vfprintf
#define _vstprintf  vsprintf
#define _vsntprintf vsnprintf
#define _tscanf     scanf
#define _ftscanf    fscanf
#define _stscanf    sscanf
#define _fgettc     fgetc
#define _fgettchar  fgetchar
#define _fgetts     fgets
#define _fputtc     fputc
#define _fputtchar  fputchar
#define _fputts     fputs
#define _tfdopen    fdopen
#define _tfopen     fopen
#define _tfreopen   freopen
#define _tfsopen    fsopen
#define _tgetenv    getenv
#define _tputenv    putenv
#define _tsearchenv searchenv
#define _tsystem	system
#define _tmakepath  makepath
#define _tsplitpath splitpath
#define _tfullpath  fullpath
#define _gettc      getc
#define _getts      gets
#define _puttc      putc
#define _puttchar	putchar
#define _putts      puts
#define _ungettc    ungetc
#define _tcstod     strtod
#define _tcstol     strtol
#define _tcstoul    strtoul
#define _itot       itoa
#define _ltot       ltoa
#define _ultot      ultoa
#define _ttoi       atoi
#define _ttol       atol
#define _tcscat     strcat
#define _tcschr     strchr
#define _tcscmp     strcmp
#define _tcscpy     strcpy
#define _tcscspn    strcspn
#define _tcslen     strlen
#define _tcsncat    strncat
#define _tcsncmp    strncmp
#define _tcsncpy    strncpy
#define _tcspbrk    strpbrk
#define _tcsrchr    strrchr
#define _tcsspn     strspn
#define _tcsstr     strstr
#define _tcstok     strtok
#define _tcsdup     strdup
#define _tcsicmp    stricmp
#define _tcsnicmp   strnicmp
#define _tcsnset    strnset
#define _tcsrev     strrev
#define _tcsset     strset
#define _tcslwr     strlwr
#define _tcsupr     strupr
#define _tcsxfrm    strxfrm
#define _tcscoll    strcoll
#define _tcsicoll   stricoll
#define _istalpha   isalpha
#define _istupper   isupper
#define _istlower   islower
#define _istdigit   isdigit
#define _istxdigit  isxdigit
#define _istspace   isspace
#define _istpunct   ispunct
#define _istalnum   isalnum
#define _istprint   isprint
#define _istgraph   isgraph
#define _istcntrl   iscntrl
#define _istascii   isascii
#define _totupper   toupper
#define _totlower   tolower
#define _tasctime   asctime
#define _tctime     ctime
#define _tstrdate   strdate
#define _tstrtime   strtime
#define _tutime     utime
#define _tcsftime   strftime

#define _tchmod     chmod
#define _tcreat     creat
#define _tfindfirst findfirst
#define _tfindnext  findnext
#define _tmktemp    mktemp
#define _topen      open
#define _taccess    access
#define _tremove    remove
#define _trename    rename
#define _tsopen     sopen
#define _tsetlocale setlocale
#define _tunlink    unlink
#define _tfinddata_t	finddata_t
#define _tchdir     chdir
#define _tgetcwd    getcwd
#define _tgetdcwd   getdcwd
#define _tmkdir     mkdir
#define _trmdir     rmdir
#define _tstat      stat

#define _stat		stat // for struct _stat

/* dirent structures and functions */
#define _tdirent    dirent
#define _TDIR       DIR
#define _topendir   opendir
#define _tclosedir  closedir
#define _treaddir   readdir
#define _trewinddir rewinddir
#define _ttelldir   telldir
#define _tseekdir   seekdir

#endif  /* Not FAKE_TCHAR_H */
