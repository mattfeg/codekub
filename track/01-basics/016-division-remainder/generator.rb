require "json"

rng = Random.new(16)
max = 10**9
pairs = [[17, 5], [4, 9], [0, 1], [10, 10], [max, 1], [max, max], [max, 7], [1, max]]
6.times { pairs << [rng.rand(0..max), rng.rand(1..1000)] }
puts JSON.generate(pairs.map { |a, b| "#{a} #{b}\n" })
