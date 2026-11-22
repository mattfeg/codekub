a, b, c, d = STDIN.read.split.map(&:to_i)
puts "X = #{a + b * c - d}"
puts "Y = #{(a + b) * (c - d)}"
puts "Z = #{(a * b + c) % d}"
puts "W = #{a + b / c * d}"
