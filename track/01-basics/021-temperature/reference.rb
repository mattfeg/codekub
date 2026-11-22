c = STDIN.gets.to_f
puts format("Fahrenheit: %.2f", c * 9 / 5 + 32)
puts format("Kelvin: %.2f", c + 273.15)
