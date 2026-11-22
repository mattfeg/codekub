require "json"

rng = Random.new(14)
ns = [7, -3, 0, 1, 10_000, -10_000]
6.times { ns << rng.rand(-10_000..10_000) }
puts JSON.generate(ns.map { |n| "#{n}\n" })
