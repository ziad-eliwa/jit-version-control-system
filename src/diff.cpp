#include "diff.h"
#include <algorithm>

Vector<std::string> diff(const Vector<std::string> &file1,
                         const Vector<std::string> &file2) {
  int n = file1.size(), m = file2.size();
  int max = n + m;
  Vector<int> del(2 * max + 1, 0);
  Vector<Vector<int>> trace;
  int x, y;
  int found = -1;
  for (int d = 0; d <= max; d++) {
    for (int k = max - d; k <= max + d; k += 2) {
      if (k == max - d || (k != max + d && del[k - 1] < del[k + 1]))
        x = del[k + 1];
      else
        x = del[k - 1] + 1;
      y = x + max - k;
      while (x < n && y < m && file1[x] == file2[y])
        x++, y++;
      del[k] = x;
      if (x >= n && y >= m) {
        found = k;
        break;
      }
    }
    trace.push_back(del);
    if (found != -1)
      break;
  }
  Vector<std::string> lines;
  int prev_x, prev_y, k = found, prev_k;
  for (int d = trace.size() - 1; d > 0; d--) {
    x = trace[d][k];
    if (k == max - d || (k != max + d && trace[d][k - 1] < trace[d][k + 1]))
      prev_k = k + 1;
    else
      prev_k = k - 1;
    prev_x = trace[d][prev_k];
    y = x + max - k, prev_y = prev_x + max - prev_k;
    while (x > prev_x && y > prev_y) {
      x--, y--;
      lines.push_back(" " + file1[x]);
    }
    if (x > prev_x) {
      x--;
      lines.push_back("-" + file1[x]);
    } else if (y > prev_y) {
      y--;
      lines.push_back("+" + file2[y]);
    }
    k = prev_k;
  }
  x--, y--;
  while (x == y && x >= 0) {
    lines.push_back(" " + file1[x]);
    x--, y--;
  }
  std::reverse(lines.begin(), lines.end());
  return lines;
}
