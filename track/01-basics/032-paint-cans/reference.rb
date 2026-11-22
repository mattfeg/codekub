length, width, height, coats, coverage = STDIN.read.split.map(&:to_i)
area = 2 * (length + width) * height * coats
puts "Área total: #{area} m²"
puts "Latas: #{(area + coverage - 1) / coverage}"
