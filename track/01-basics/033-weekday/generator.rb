require "json"

rng = Random.new(33)
max = 10**9
cases = [[3, 10], [1, -3], [0, 0], [6, 1], [0, -1], [0, -7], [6, max], [0, -max]]
7.times { cases << [rng.rand(0..6), rng.rand(-max..max)] }
puts JSON.generate(cases.map { |d, n| "#{d} #{n}\n" })
