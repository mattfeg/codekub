require "json"

rng = Random.new(12)
names = ["Carla Mendes", "Diego", "Eduarda Lima Prado", "Fábio", "Gabriela Antunes", "Heitor Castro"]
cities = ["Manaus", "Porto Alegre", "Belo Horizonte", "Natal", "Rio de Janeiro", "Cuiabá", "Florianópolis"]
cases = [["Ana Souza", 28, "Recife"], ["Bruno", 0, "São José dos Campos"], ["Idosa Feliz", 120, "Salvador"]]
8.times { cases << [names[rng.rand(names.size)], rng.rand(0..120), cities[rng.rand(cities.size)]] }
puts JSON.generate(cases.map { |n, a, c| "#{n}\n#{a}\n#{c}\n" })
