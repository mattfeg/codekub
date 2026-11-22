require "json"

rng = Random.new(11)
max = 10**9
pairs = [[3, 7], [-1, -1], [0, 5], [max, -max], [42, 0]]
7.times { pairs << [rng.rand(-max..max), rng.rand(-max..max)] }
puts JSON.generate(pairs.map { |a, b| "#{a} #{b}\n" })
