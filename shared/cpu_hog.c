#include <stdlib.h>

#ifdef USE_GMP
#include <gmp.h>
#endif

long *side_effect;

// cpuHog tries to waste about 1ms of CPU time.
long cpuHog() {
	// call plain malloc, which gives us an invocation statically built in
	// to the program
	long *p = malloc(1000*sizeof(long));
	side_effect = p;
	free(p);

	long sum = 0;
	#ifdef USE_GMP
	// Do some GMP stuff to get memory allocation from a third-party shared
	// library. We use p to make sure all the operations in this function
	// are intertwined and hopefully don't get optimized away
	mpz_t foo;
	mpz_init_set_ui(foo, (unsigned long int) p);
	mpz_pow_ui(foo, foo, 42);
	sum = mpz_odd_p(foo);
	#endif

	for (long i = 0; i < 1000000; i++) {
		sum += i % 10;
	}

	#ifdef USE_GMP
	mpz_clear(foo);
	#endif

	return sum;
}