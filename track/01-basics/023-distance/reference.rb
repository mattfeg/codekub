x1, y1, x2, y2 = STDIN.read.split.map(&:to_f)
puts format("%.4f", Math.sqrt((x2 - x1)**2 + (y2 - y1)**2))
