/* Parses the raw data file produced by LTSpice and prints it as text.
 * 
 * R.Veen, 4/2008.
 * License: public domain.
 */

#include <stdlib.h>
#include <stdio.h>
#include <string.h>

int n=0, m=0;

void parse_text(int noText)
{
	int i;
	char buf[256], *p, *q, *r;
	char lab[256][128];
	char unit[256][128];
	
	while ( fgets(buf,256,stdin) != NULL ) {
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
				fgets(buf,256,stdin);
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
	
	// printf("data_x %d\n",n);
	// printf("data_y %d\n",m);
	// printf("labels\n");

	if (!noText) {
	for (i=0; i<n; i++) {
	    printf("\"%s\"",lab[i]);
	    if (i<n-1)
	    	printf(", ");
	}
	putchar('\n');
	}
	/*
	printf("quantity\n");
	for (i=0; i<n; i++) 
	    printf("  %s\n",unit[i]);
	    */
}

void parse_binary()
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
		i = fread(u.b, 8, 1, stdin);
	        if (!i) break;
	        
	        printf("%f, ",u.d);
	    } else {
		    i = fread(u.b, 4, 1, stdin);
		    if (!i) break;

		    if (tf>=n-1) {
		      tf = -1;
		      printf("%f\n",u.f);
		    }
		    else {
		    	printf("%f, ",u.f);
		    }
	    }
	    
	    tf++;
	}
}

int main(int argc, char**argv)
{
	int noText = 1;

	if (argc>1) {
		if (!strcmp("-t",argv[1]))
			noText = 0;
		else {
		    puts("usage: lt2csv < file.raw");
		    puts("       lt2csv -t < file.raw [header only]");
		    exit(0);
		}
	}

    parse_text(noText);

    parse_binary();
}
