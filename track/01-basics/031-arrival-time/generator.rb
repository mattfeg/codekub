require "json"

rng = Random.new(31)
cases = [[22, 30, 525], [8, 0, 90], [0, 0, 0], [23, 59, 1], [12, 0, 1440], [23, 59, 1_000_000]]
8.times { cases << [rng.rand(0..23), rng.rand(0..59), rng.rand(0..1_000_000)] }
puts JSON.generate(cases.map { |h, m, d| "#{h} #{m}\n#{d}\n" })
