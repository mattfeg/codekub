n, m, a = STDIN.read.split.map(&:to_i)
puts ((n + a - 1) / a) * ((m + a - 1) / a)
