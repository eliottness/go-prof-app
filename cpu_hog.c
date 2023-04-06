#include <stdlib.h>

long *side_effect;

// We give a pointer to this function via the go call
// it prevent further optimization by the compiler
long cpuHogOp(long a, long b) { return a + b % 10; }

// cpuHog tries to waste about 1ms of CPU time.
long cpuHog(long (*op)(long, long)) {
  // call plain malloc, which gives us an invocation statically built in
  // to the program
  long *p = malloc(1000 * sizeof(long));
  side_effect = p;
  free(p);

  long sum = 0;

  // Prime number to prevent easy loop unwrapping
  for (long i = 0; i < 999983; i++) {
    sum = op(sum, i);

    // Unpredicable branch
    if (sum % 10 != 0)
        return 0;
  }

  return sum;
}