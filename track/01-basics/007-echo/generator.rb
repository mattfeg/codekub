require "json"

rng = Random.new(7)
words = %w[código programa teste juiz entrada saída linha texto Ruby Python Go laço
           variável função ação maçã pão você café 42 2026 ok! sim? não. eco, olá]
inputs = ["Olá, codekub!", "Aprender a programar é praticar todo dia", "a", "x" * 100,
          "1 2 3 4 5", "Pontuação: vírgula, ponto; e exclamação!"]
8.times do
  phrase = Array.new(rng.rand(2..12)) { words[rng.rand(words.size)] }.join(" ")
  inputs << phrase[0, 100].strip
end
puts JSON.generate(inputs.map { |s| s + "\n" })
