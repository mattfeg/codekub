require "json"

rng = Random.new(19)
cases = [[2, 3, 4, 5], [10, 7, 3, 4], [1, 1, 1, 1], [1000, 1000, 1000, 1000], [5, 999, 1000, 7], [1, 1000, 1, 1000]]
8.times { cases << Array.new(4) { rng.rand(1..1000) } }
puts JSON.generate(cases.map { |c| c.join(" ") + "\n" })
