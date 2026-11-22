require "json"

rng = Random.new(9)
max = 10**9
pairs = [[3, 5], [-700, 250], [0, 0], [max, max], [-max, -max], [max, -max], [1, -1]]
4.times { pairs << [rng.rand(-100..100), rng.rand(-100..100)] }
4.times { pairs << [rng.rand(-max..max), rng.rand(-max..max)] }
puts JSON.generate(pairs.map { |a, b| "#{a} #{b}\n" })
