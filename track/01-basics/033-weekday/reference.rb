today, n = STDIN.read.split.map(&:to_i)
puts "dia: #{(today + n) % 7}"
