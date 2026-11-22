require "json"

rng = Random.new(29)
vs = [388, 76, 0, 1, 200, 399, 1_000_000, 999_999]
6.times { vs << rng.rand(0..1_000_000) }
puts JSON.generate(vs.map { |v| "#{v}\n" })
