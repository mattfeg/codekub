items = [
  ["Café", 4.50],
  ["Pão de queijo", 6.00],
  ["Bolo de milho", 12.90],
  ["Suco natural", 8.75]
]
border = "+" + "-" * 15 + "+" + "-" * 9 + "+"

puts border
puts "| #{'Produto'.ljust(13)} | #{'Preço'.rjust(7)} |"
puts border
items.each do |name, price|
  puts "| #{name.ljust(13)} | #{format('%.2f', price).rjust(7)} |"
end
puts border
