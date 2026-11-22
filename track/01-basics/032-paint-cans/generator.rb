require "json"

rng = Random.new(32)
cases = [[4, 3, 3, 2, 36], [5, 5, 2, 1, 40], [1, 1, 1, 1, 1], [1, 1, 1, 1, 1000], [100, 100, 100, 5, 1], [100, 100, 100, 5, 1000]]
8.times { cases << [rng.rand(1..100), rng.rand(1..100), rng.rand(1..100), rng.rand(1..5), rng.rand(1..1000)] }
puts JSON.generate(cases.map { |c, w, h, d, r| "#{c} #{w} #{h}\n#{d} #{r}\n" })
