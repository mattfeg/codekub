require "json"

rng = Random.new(17)
ns = [407, 999, 100, 101, 110, 909, 550]
6.times { ns << rng.rand(100..999) }
puts JSON.generate(ns.map { |n| "#{n}\n" })
