require "json"

rng = Random.new(15)
cases = [[70, 80, 90], [100, 100, 99], [0, 0, 0], [100, 100, 100], [0, 0, 2], [1, 1, 0]]
7.times { cases << Array.new(3) { rng.rand(0..100) } }
puts JSON.generate(cases.map { |c| c.join(" ") + "\n" })
