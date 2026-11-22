require "json"

rng = Random.new(8)
first = %w[Ana João Maria Pedro Luíza Rafael Beatriz Tiago Camila Ícaro Zé Yasmin]
last = %w[Souza Silva Oliveira Santos Ferreira Lima Araújo Melo Gonçalves]
inputs = ["Ana", "João Pedro", "Z", "Maria Clara dos Santos Ferreira"]
8.times do
  parts = [first[rng.rand(first.size)]] + Array.new(rng.rand(0..2)) { last[rng.rand(last.size)] }
  inputs << parts.join(" ")
end
puts JSON.generate(inputs.map { |s| s + "\n" })
