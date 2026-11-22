require "json"

rng = Random.new(10)
max = 10**6
cases = [[100, 250, 75], [0, 0, 1], [0, 0, 0], [max, max, max], [1, 2, 3]]
8.times { cases << Array.new(3) { rng.rand(0..max) } }
puts JSON.generate(cases.map { |c| c.join("\n") + "\n" })
