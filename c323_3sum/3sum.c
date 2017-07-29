#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <string.h>

typedef struct {
  long long int *nums;
  size_t len;
} nums_s;

int cmp_lli(const void *lp, const void *rp) {
    long long int l = *(const long long int *)lp;
    long long int r = *(const long long int *)rp;

    if (l < r) return -1;
    else if (l > r) return 1;
    return 0;
}

int cmp_3lli(const void *lp, const void *rp) {
    const long long int *l = *(const long long int **)lp;
    const long long int *r = *(const long long int **)rp;

    if (l[0] < r[0]) return -1;
    else if (l[0] > r[0]) return 1;
    else {
      if (l[1] < r[1]) return -1;
      else if (l[1] > r[1]) return 1;
      else {
        if (l[2] < r[2]) return -1;
        else if (l[2] > r[2]) return 1;
      }
    }
    return 0;
}

nums_s *nums_new_from_line(char const *line) {
  size_t len = 0;
  size_t cap = 16;
  long long int *nums = malloc(cap * sizeof(long long int));

  char const *lp = line;
  char *endptr;

  while (*lp != '\n') {
    errno = 0;
    long long int n = strtoll(lp, &endptr, 10);
    if (errno) {
      fprintf(stderr, "%s\n", strerror(errno));
      exit(1);
    }
    if (lp == endptr) {
      fprintf(stderr, "can't parse %s\n", lp);
      exit(1);
    }

    lp = endptr;

    if (len == cap) {
      cap = cap * 3 / 2;
      nums = realloc(nums, cap * sizeof(long long int));
    }
    nums[len] = n;
    len++;
  }

  qsort(nums, len, sizeof(long long int), &cmp_lli);

  nums_s *ns = malloc(sizeof(nums_s));
  ns->nums = nums;
  ns->len = len;

  return ns;
}

void nums_free(nums_s *n) {
  free(n->nums);
  free(n);
}

int main() {
  size_t line_cap = 0;
  char *line = NULL;

  while (1) {
    ssize_t line_len = getline(&line, &line_cap, stdin);
    if (line_len == -1) {
      if (errno != 0) {
        fprintf(stderr, "%s\n", strerror(errno));
        exit(1);
      }
      break;
    } else if (line_len == 1) {
      continue;
    }

    nums_s *nums = nums_new_from_line(line);

    size_t len = 0;
    size_t cap = 16;
    long long int **answers = malloc(cap * sizeof(long long int*));

    for (size_t i = 0; i < nums->len; i++) {
      for (size_t j = i + 1; j < nums->len; j++) {
        for (size_t k = j + 1; k < nums->len; k++) {
          long long int sum = nums->nums[i] + nums->nums[j] + nums->nums[k];
          if (sum == 0) {
            if (len == cap) {
              cap = cap * 3 / 2;
              answers = realloc(answers, cap * sizeof(long long int*));
            }
            long long int *a = malloc(3 * sizeof(long long int));
            a[0] = nums->nums[i];
            a[1] = nums->nums[j];
            a[2] = nums->nums[k];
            answers[len] = a;
            len++;
          }
        }
      }
    }

    nums_free(nums);

    if (len > 0) {
      qsort(answers, len, sizeof(long long int *), &cmp_3lli);

      long long int *last = answers[0];
      printf("%lli %lli %lli\n", last[0], last[1], last[2]);

      for (size_t i = 1; i < len; i++) {
        long long int *curr = answers[i];
        if (memcmp(last, curr, 3 * sizeof(long long int))) {
          printf("%lli %lli %lli\n", curr[0], curr[1], curr[2]);
          last = curr;
        }
      }
    }

    for (size_t i = 0; i < len; i++) {
      free(answers[i]);
    }
    free(answers);

    printf("\n");
  }

  free(line);
}
