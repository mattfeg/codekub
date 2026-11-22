require "json"

def safe?(v, d)
  return false if format("%.#{d}f", v) =~ /\A-0\.?0*\z/
  frac = (v.abs * 10**d) % 1
  (frac - 0.5).abs > 1e-6
end

dist = ->(p) { Math.sqrt((p[2] - p[0])**2 + (p[3] - p[1])**2) }
rng = Random.new(23)
cases = [[0.0, 0.0, 3.0, 4.0], [1.5, -2.0, -0.5, 3.0], [7.0, 7.0, 7.0, 7.0], [-1000.0, -1000.0, 1000.0, 1000.0]]
while cases.size < 13
  p = Array.new(4) { rng.rand(-10_000..10_000) / 10.0 }
  cases << p if safe?(dist.call(p), 4)
end
abort "tie in fixed case" unless cases.all? { |p| safe?(dist.call(p), 4) }
puts JSON.generate(cases.map { |p| format("%.1f %.1f\n%.1f %.1f\n", *p) })
