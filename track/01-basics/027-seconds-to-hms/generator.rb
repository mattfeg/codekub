require "json"

rng = Random.new(27)
ns = [3725, 0, 59, 60, 3599, 3600, 86_399, 86_400, 10_000_000]
5.times { ns << rng.rand(0..10_000_000) }
puts JSON.generate(ns.map { |n| "#{n}\n" })
