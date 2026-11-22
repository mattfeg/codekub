years, months, days = STDIN.read.split.map(&:to_i)
puts "dias vividos: #{years * 365 + months * 30 + days}"
