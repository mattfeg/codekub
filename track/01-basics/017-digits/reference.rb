n = STDIN.gets.to_i
hundreds = n / 100
tens = n / 10 % 10
units = n % 10
puts "centena: #{hundreds}"
puts "dezena: #{tens}"
puts "unidade: #{units}"
puts "soma: #{hundreds + tens + units}"
