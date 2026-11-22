require "json"

rng = Random.new(18)
ns = [3, 1000, 0, 1, 2, 10, 999]
5.times { ns << rng.rand(0..1000) }
puts JSON.generate(ns.map { |n| "#{n}\n" })
