/* Parses the raw data file produced by LTSpice and prints it as text.
 * 
 * R.Veen, 4/2008.
 * License: public domain.
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

int n=0, m=0;

void parse_text(int print, FILE *in)
{
	int i;
	char buf[256], *p, *q, *r;
	char lab[256][128];
	char unit[256][128];
	
	while ( fgets(buf,256,in) != NULL ) {
		if ( !strncmp("Binary:",buf,7))
			break;
		
		if (!strncmp("No. Variables:",buf,14)) {
			n = atoi(buf+15);
		}
		else if (!strncmp("No. Points:",buf,11)) {
			m = atoi(buf+12);
		}
		else if (!strncmp("Variables:",buf,10)) {
			for (i=0; i<n; i++) {
				fgets(buf,256,in);
				p = buf+1;
				while (*p && !isblank(*p))
					p++;
				while (*p && isblank(*p))
					p++;
				
				q = p+1;
				while (*q && !isblank(*q))
					q++;
				*q=0;
				q++;
				while (*q && isblank(*q))
					q++;
				
				r = q;
				while (*r && !isblank(*r))
				    r++;
				*--r = 0;
				strncpy(lab[i],p,128);
				strncpy(unit[i],q,128);
			}
		}
	}
	
	if (!print)
		return;

	printf("data_x %d\n",n);
	printf("data_y %d\n",m);
	printf("columns\n");

	for (i=0; i<n; i++) {
	    printf("  '%s'\n",lab[i]);
	}
}

void parse_binary(FILE *in)
{
	char c, tf=0, i;
		
	union {
		float f;
		char b[8];
		double d;
	} u;
	
	while ( 1 ) {
		
            /* first float is 8 bytes (time),
             * rest is 4 byte.
             */

	    if (!tf) { 
		i = fread(u.b, 8, 1, in);
	        if (!i) break;
	        
	        printf("%.12f, ",u.d);
	    } else {
		    i = fread(u.b, 4, 1, in);
		    if (!i) break;

		    if (tf>=n-1) {
		      tf = -1;
		      printf("%.12g\n",u.f);
		    }
		    else {
		    	printf("%.12g, ",u.f);
		    }
	    }
	    
	    tf++;
	}
}

int main(int argc, char**argv)
{
	FILE *f;
	int header = 0, ix=1;

	while (argc>1) {
		if (!strcmp("-t",argv[ix])) {
			header = 1;
			argc--;
			ix++;
			continue;
		}

		if (!strcmp("-h",argv[ix])) {
			puts("usage: lt2csv [<] file.raw");
			puts("       lt2csv -t [<] file.raw [prints header only]");
			exit(0);
		}

		f = fopen(argv[ix],"r");

		if (header) {
		    fseek(f,0L,SEEK_END);
		    printf("bytes %ld\n",ftell(f));
		    rewind(f);
		}

	    break;
	}

	if (f==NULL) {
		f = stdin;
	}

	parse_text(header,f);

	if (!header)
        parse_binary(f);
}
