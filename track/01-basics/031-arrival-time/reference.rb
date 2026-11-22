h, m, duration = STDIN.read.split.map(&:to_i)
total = h * 60 + m + duration
puts format("chegada: %02d:%02d", total / 60 % 24, total % 60)
puts "dias depois: #{total / 1440}"
