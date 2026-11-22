require "json"

rng = Random.new(25)
xs = ["2.5", "-2.5", "3.141", "-3.141", "0.5", "-0.5", "-0.3", "7", "-1000", "999.999", "0.499"]
4.times { xs << format("%.3f", rng.rand(-1_000_000..1_000_000) / 1000.0) }
puts JSON.generate(xs.map { |x| x + "\n" })
