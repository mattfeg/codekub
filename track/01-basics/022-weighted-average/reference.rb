a, b, c = STDIN.read.split.map(&:to_f)
puts format("Média final: %.1f", (3 * a + 3 * b + 4 * c) / 10)
