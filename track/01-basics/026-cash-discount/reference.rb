price, percent = STDIN.read.split
price = price.to_f
discount = price * percent.to_i / 100
puts format("Desconto: R$ %.2f", discount)
puts format("Total: R$ %.2f", price - discount)
