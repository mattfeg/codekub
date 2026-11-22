require "json"

rng = Random.new(13)
max = 10**9
ns = [10, 0, -1, 1, max, -max]
6.times { ns << rng.rand(-max..max) }
puts JSON.generate(ns.map { |n| "#{n}\n" })
