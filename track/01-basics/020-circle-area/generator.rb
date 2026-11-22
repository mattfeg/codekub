require "json"

def safe?(v, d)
  return false if format("%.#{d}f", v) =~ /\A-0\.?0*\z/
  frac = (v.abs * 10**d) % 1
  (frac - 0.5).abs > 1e-6
end

rng = Random.new(20)
rs = [2.0, 0.5, 0.01, 10_000.0, 1.0]
while rs.size < 13
  r = rng.rand(1..1_000_000) / 100.0
  rs << r if safe?(3.14159 * r * r, 4)
end
abort "tie in fixed case" unless rs.all? { |r| safe?(3.14159 * r * r, 4) }
puts JSON.generate(rs.map { |r| format("%.2f\n", r) })
