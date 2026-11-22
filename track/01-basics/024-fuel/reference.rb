km = STDIN.gets.to_i
liters = STDIN.gets.to_f
price = STDIN.gets.to_f
puts format("Consumo: %.3f km/l", km / liters)
puts format("Custo por km: R$ %.2f", liters * price / km)
