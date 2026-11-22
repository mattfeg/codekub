name = STDIN.gets.chomp
age = STDIN.gets.to_i
city = STDIN.gets.chomp
puts "Nome: #{name}"
puts "Idade: #{age} anos"
puts "Cidade: #{city}"
puts "Em 10 anos, #{name} terá #{age + 10} anos."
