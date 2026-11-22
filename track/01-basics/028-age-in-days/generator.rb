require "json"

rng = Random.new(28)
cases = [[20, 3, 15], [0, 0, 0], [0, 0, 29], [0, 11, 29], [120, 11, 29], [1, 0, 0]]
7.times { cases << [rng.rand(0..120), rng.rand(0..11), rng.rand(0..29)] }
puts JSON.generate(cases.map { |c| c.join(" ") + "\n" })
